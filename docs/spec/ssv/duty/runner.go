package duty

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// DutyExecutionSlotTimeout is the timeout for pre or post consensus signature collection.
const DutyExecutionSlotTimeout spec.Slot = 32

// Runner is manages the execution of a duty from start to finish, it can only execute 1 duty at a time.
// Prev duty must finish before the next one can start.
type Runner struct {
	BeaconRoleType beacon.RoleType
	BeaconNetwork  ssv.BeaconNetwork
	Share          *types.Share
	// State holds all relevant params for a full duty execution (consensus & post consensus)
	State *State
	// CurrentDuty is the current executing duty, changes once StartNewDuty is called
	CurrentDuty    *beacon.Duty
	QBFTController *qbft.Controller
	storage        ssv.Storage
	valCheck       qbft.ProposedValueCheck
}

func NewDutyRunner(
	beaconRoleType beacon.RoleType,
	beaconNetwork ssv.BeaconNetwork,
	share *types.Share,
	qbftController *qbft.Controller,
	storage ssv.Storage,
	valCheck qbft.ProposedValueCheck,
) *Runner {
	return &Runner{
		BeaconRoleType: beaconRoleType,
		BeaconNetwork:  beaconNetwork,
		Share:          share,
		QBFTController: qbftController,
		storage:        storage,
		valCheck:       valCheck,
	}
}

func (dr *Runner) StartNewDuty(duty *beacon.Duty) error {
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
func (dr *Runner) CanStartNewDuty(duty *beacon.Duty) error {
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

// GetRoot returns the root used for signing and verification
func (dr *Runner) GetRoot() ([]byte, error) {
	marshaledRoot, err := dr.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Encode returns the encoded struct in bytes or error
func (dr *Runner) Encode() ([]byte, error) {
	return json.Marshal(dr)
}

// Decode returns error if decoding failed
func (dr *Runner) Decode(data []byte) error {
	return json.Unmarshal(data, &dr)
}

func (dr *Runner) validatePartialSigMsg(msg *ssv.SignedPartialSignatureMessage, container *PartialSigContainer) error {
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

// postConsensusSigTimeout returns true if collecting sigs timed out
func (dr *Runner) partialSigCollectionTimeout(container *PartialSigContainer, currentSlot spec.Slot) bool {
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
