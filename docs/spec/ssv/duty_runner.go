package ssv

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// DutyExecutionSlotTimeout is the timeout for pre or post consensus signature collection.
const DutyExecutionSlotTimeout spec.Slot = 32

// DutyRunner is manages the execution of a duty from start to finish, it can only execute 1 duty at a time.
// Prev duty must finish before the next one can start.
type DutyRunner struct {
	BeaconRoleType beacon.RoleType
	BeaconNetwork  BeaconNetwork
	Share          *types.Share
	// State holds all relevant params for a full duty execution (consensus & post consensus)
	State *RunnerState
	// CurrentDuty is the current executing duty, changes once StartNewDuty is called
	CurrentDuty    *beacon.Duty
	QBFTController *qbft.Controller
	storage        Storage
	valCheck       qbft.ProposedValueCheck
}

func NewDutyRunner(
	beaconRoleType beacon.RoleType,
	beaconNetwork BeaconNetwork,
	share *types.Share,
	qbftController *qbft.Controller,
	storage Storage,
	valCheck qbft.ProposedValueCheck,
) *DutyRunner {
	return &DutyRunner{
		BeaconRoleType: beaconRoleType,
		BeaconNetwork:  beaconNetwork,
		Share:          share,
		QBFTController: qbftController,
		storage:        storage,
		valCheck:       valCheck,
	}
}

func (dr *DutyRunner) StartNewDuty(duty *beacon.Duty) error {
	if err := dr.CanStartNewDuty(duty); err != nil {
		return err
	}
	dr.CurrentDuty = duty
	dr.State = NewDutyExecutionState(dr.Share.Quorum)
	return nil
}

// CanStartNewDuty returns nil if:
// - no running instance exists or
// - pre consensus timeout
// - a QBFT instance Decided and all post consensus sigs collectd or
// - a QBFT instance Decided and post consensus timeout
// else returns an error
func (dr *DutyRunner) CanStartNewDuty(duty *beacon.Duty) error {
	if dr.State == nil || dr.State.IsFinished() {
		return nil
	}

	// if running instance hasn't decided yet we always return error
	if dr.State.RunningInstance != nil {
		// check consensus decided
		if decided, _ := dr.State.RunningInstance.IsDecided(); !decided {
			return errors.New("consensus on duty is running")
		}
	}

	// check pre consensus signature collection timeout
	switch dr.BeaconRoleType {
	case beacon.RoleTypeProposer:
		if dr.randaoSigTimeout(duty.Slot) {
			return nil
		}
		if !dr.State.RandaoPartialSig.HasQuorum() {
			return errors.New("randao consensus sig collection is running")
		}
	}

	// check if completed post consensus sigs or timedout
	if !dr.State.PostConsensusPartialSig.HasQuorum() && !dr.postConsensusSigTimeout(duty.Slot) {
		return errors.New("post consensus sig collection is running")
	}
	return nil
}

// StartNewConsensusInstance starts a new QBFT instance for value
func (dr *DutyRunner) StartNewConsensusInstance(value []byte) error {
	if len(value) == 0 {
		return errors.New("new instance value invalid")
	}
	if err := dr.QBFTController.StartNewInstance(value); err != nil {
		return errors.Wrap(err, "could not start new QBFT instance")
	}
	newInstance := dr.QBFTController.InstanceForHeight(dr.QBFTController.Height)
	if newInstance == nil {
		return errors.New("could not find newly created QBFT instance")
	}

	dr.State.RunningInstance = newInstance
	return nil
}

// ProcessRandaoMessage process randao msg, returns true if it has quorum for partial signatures.
// returns true only once (first time quorum achieved)
func (dr *DutyRunner) ProcessRandaoMessage(msg *SignedPartialSignatureMessage) (bool, error) {
	if err := dr.canProcessRandaoMsg(msg); err != nil {
		return false, errors.Wrap(err, "can't process randao message")
	}

	prevQuorum := dr.State.RandaoPartialSig.HasQuorum()

	if err := dr.State.RandaoPartialSig.AddSignature(msg.Message); err != nil {
		return false, errors.Wrap(err, "could not add partial randao signature")
	}

	if prevQuorum {
		return false, nil
	}

	return dr.State.RandaoPartialSig.HasQuorum(), nil
}

func (dr *DutyRunner) ProcessConsensusMessage(msg *qbft.SignedMessage) (decided bool, decidedValue *types.ConsensusData, err error) {
	decided, decidedValueByts, err := dr.QBFTController.ProcessMsg(msg)
	if err != nil {
		return false, nil, errors.Wrap(err, "failed to process consensus msg")
	}

	/**
	Decided returns true only once so if it is true it must be for the current running instance
	*/
	if !decided {
		return false, nil, nil
	}

	decidedValue = &types.ConsensusData{}
	if err := decidedValue.Decode(decidedValueByts); err != nil {
		return true, nil, errors.Wrap(err, "failed to parse decided value to ConsensusData")
	}

	if err := dr.validateDecidedConsensusData(decidedValue); err != nil {
		return true, nil, errors.Wrap(err, "decided ConsensusData invalid")
	}

	return true, decidedValue, nil
}

// ProcessPostConsensusMessage process post consensus msg, returns true if it has quorum for partial signatures.
// returns true only once (first time quorum achieved)
func (dr *DutyRunner) ProcessPostConsensusMessage(msg *SignedPartialSignatureMessage) (bool, error) {
	if err := dr.canProcessPostConsensusMsg(msg); err != nil {
		return false, errors.Wrap(err, "can't process post consensus message")
	}

	prevQuorum := dr.State.PostConsensusPartialSig.HasQuorum()

	if err := dr.State.PostConsensusPartialSig.AddSignature(msg.Message); err != nil {
		return false, errors.Wrap(err, "could not add partial post consensus signature")
	}

	if prevQuorum {
		return false, nil
	}

	dr.State.SetFinished()
	return dr.State.PostConsensusPartialSig.HasQuorum(), nil
}

// SignDutyPostConsensus sets the Decided duty and partially signs the Decided data, returns a PartialSignatureMessage to be broadcasted or error
func (dr *DutyRunner) SignDutyPostConsensus(decidedValue *types.ConsensusData, signer types.KeyManager) (*PartialSignatureMessage, error) {
	ret := &PartialSignatureMessage{
		Type:    PostConsensusPartialSig,
		Signers: []types.OperatorID{dr.Share.OperatorID},
	}

	switch dr.BeaconRoleType {
	case beacon.RoleTypeAttester:
		signedAttestation, r, err := signer.SignAttestation(decidedValue.AttestationData, decidedValue.Duty, dr.Share.SharePubKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign attestation")
		}

		dr.State.DecidedValue = decidedValue
		dr.State.SignedAttestation = signedAttestation
		dr.State.PostConsensusPartialSig.SigRoot = ensureRoot(r)

		ret.SigningRoot = dr.State.PostConsensusPartialSig.SigRoot
		ret.PartialSignature = dr.State.SignedAttestation.Signature[:]

		return ret, nil
	default:
		return nil, errors.Errorf("unknown duty %s", decidedValue.Duty.Type.String())
	}
}

// GetRoot returns the root used for signing and verification
func (dr *DutyRunner) GetRoot() ([]byte, error) {
	marshaledRoot, err := dr.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Encode returns the encoded struct in bytes or error
func (dr *DutyRunner) Encode() ([]byte, error) {
	return json.Marshal(dr)
}

// Decode returns error if decoding failed
func (dr *DutyRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &dr)
}

// canProcessRandaoMsg returns true if it can process randao message, false if not
func (dr *DutyRunner) canProcessRandaoMsg(msg *SignedPartialSignatureMessage) error {
	if err := dr.validatePartialSigMsg(msg, dr.State.RandaoPartialSig); err != nil {
		return errors.Wrap(err, "randao msg invalid")
	}

	if dr.randaoSigTimeout(dr.BeaconNetwork.EstimatedCurrentSlot()) {
		return errors.New("randao sig collection timeout")
	}

	return nil
}

// canProcessPostConsensusMsg returns true if it can process post consensus message, false if not
func (dr *DutyRunner) canProcessPostConsensusMsg(msg *SignedPartialSignatureMessage) error {
	if err := dr.validatePartialSigMsg(msg, dr.State.PostConsensusPartialSig); err != nil {
		return errors.Wrap(err, "post consensus msg invalid")
	}

	if dr.postConsensusSigTimeout(dr.BeaconNetwork.EstimatedCurrentSlot()) {
		return errors.New("post consensus sig collection timeout")
	}

	return nil
}

func (dr *DutyRunner) validatePartialSigMsg(msg *SignedPartialSignatureMessage, container *PartialSigContainer) error {
	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "SignedPartialSignatureMessage invalid")
	}

	if err := msg.GetSignature().VerifyByOperators(msg, dr.Share.DomainType, types.PartialSignatureType, dr.Share.Committee); err != nil {
		return errors.Wrap(err, "failed to verify PartialSignature")
	}

	// validate signing root equal to Decided
	if !bytes.Equal(container.SigRoot, msg.Message.SigningRoot) {
		return errors.New("partial sig signing root is wrong")
	}

	if err := dr.verifyBeaconPartialSignature(msg.Message); err != nil {
		return errors.Wrap(err, "could not verify beacon partial Signature")
	}

	return nil
}

func (dr *DutyRunner) verifyBeaconPartialSignature(msg *PartialSignatureMessage) error {
	if len(msg.Signers) != 1 {
		return errors.New("PartialSignatureMessage allows 1 signer")
	}

	signer := msg.Signers[0]
	signature := msg.PartialSignature
	root := msg.SigningRoot

	for _, n := range dr.Share.Committee {
		if n.GetID() == signer {
			pk := &bls.PublicKey{}
			if err := pk.Deserialize(n.GetPublicKey()); err != nil {
				return errors.Wrap(err, "could not deserialized pk")
			}
			sig := &bls.Sign{}
			if err := sig.Deserialize(signature); err != nil {
				return errors.Wrap(err, "could not deserialized Signature")
			}

			// protect nil root
			root = ensureRoot(root)
			// verify
			if !sig.VerifyByte(pk, root) {
				return errors.Errorf("could not verify Signature from iBFT member %d", signer)
			}
			return nil
		}
	}
	return errors.New("beacon partial Signature signer not found")
}

func (dr *DutyRunner) validateDecidedConsensusData(val *types.ConsensusData) error {
	byts, err := val.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode decided value")
	}
	if err := dr.valCheck(byts); err != nil {
		return errors.Wrap(err, "decided value is invalid")
	}

	if dr.BeaconRoleType != val.Duty.Type {
		return errors.New("decided value's duty has wrong beacon role type")
	}

	if !bytes.Equal(dr.Share.ValidatorPubKey, val.Duty.PubKey[:]) {
		return errors.New("decided value's validator pk is wrong")
	}

	return nil
}

// postConsensusSigTimeout returns true if collecting post consensus sigs timed out
func (dr *DutyRunner) postConsensusSigTimeout(currentSlot spec.Slot) bool {
	return dr.partialSigCollectionTimeout(dr.State.PostConsensusPartialSig, currentSlot)
}

// randaoSigTimeout returns true if collecting post consensus sigs timed out
func (dr *DutyRunner) randaoSigTimeout(currentSlot spec.Slot) bool {
	return dr.partialSigCollectionTimeout(dr.State.RandaoPartialSig, currentSlot)
}

// postConsensusSigTimeout returns true if collecting sigs timed out
func (dr *DutyRunner) partialSigCollectionTimeout(container *PartialSigContainer, currentSlot spec.Slot) bool {
	if dr.CurrentDuty == nil || dr.State == nil {
		return false
	}

	if container.HasQuorum() {
		return false
	}

	if dr.CurrentDuty.Slot+DutyExecutionSlotTimeout > currentSlot {
		return false
	}

	return true
}

// ensureRoot ensures that SigningRoot will have sufficient allocated memory
// otherwise we get panic from bls:
// github.com/herumi/bls-eth-go-binary/bls.(*Sign).VerifyByte:738
func ensureRoot(root []byte) []byte {
	n := len(root)
	if n == 0 {
		n = 1
	}
	tmp := make([]byte, n)
	copy(tmp[:], root[:])
	return tmp[:]
}
