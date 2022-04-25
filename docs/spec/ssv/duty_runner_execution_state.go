package ssv

import (
	"crypto/sha256"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

type PartialSigContainer struct {
	Signatures map[types.OperatorID][]byte
	SigRoot    []byte
	// Quorum is the number of min signatures needed for quorum
	Quorum uint64
}

func NewPartialSigContainer(quorum uint64) *PartialSigContainer {
	return &PartialSigContainer{
		Quorum:     quorum,
		Signatures: make(map[types.OperatorID][]byte),
	}
}

func (ps *PartialSigContainer) AddSignature(sigMsg *PartialSignatureMessage) error {
	if len(sigMsg.Signers) != 1 {
		return errors.New("PartialSignatureMessage has != 1 Signers")
	}

	if ps.Signatures[sigMsg.Signers[0]] == nil {
		ps.Signatures[sigMsg.Signers[0]] = sigMsg.PartialSignature
	}
	return nil
}

func (ps *PartialSigContainer) ReconstructSignature(validatorPubKey []byte) ([]byte, error) {
	// Reconstruct signatures
	signature, err := types.ReconstructSignatures(ps.Signatures)
	if err != nil {
		return nil, errors.Wrap(err, "failed to reconstruct signatures")
	}
	if err := types.VerifyReconstructedSignature(signature, validatorPubKey, ps.SigRoot); err != nil {
		return nil, errors.Wrap(err, "failed to verify reconstruct signature")
	}
	return signature.Serialize(), nil
}

func (ps *PartialSigContainer) HasQuorum() bool {
	return uint64(len(ps.Signatures)) >= ps.Quorum
}

// DutyExecutionState holds all the relevant progress the duty execution progress
type DutyExecutionState struct {
	// Height represents a unique consensus height for this state
	Height          qbft.Height
	RunningInstance *qbft.Instance

	ProposedValue *types.ConsensusData
	DecidedValue  *types.ConsensusData

	SignedAttestation *spec.Attestation
	SignedProposal    *spec.SignedBeaconBlock

	RandaoPartialSig        *PartialSigContainer
	PostConsensusPartialSig *PartialSigContainer

	Finished bool
}

func NewDutyExecutionState(quorum uint64, height qbft.Height) *DutyExecutionState {
	return &DutyExecutionState{
		Height:                  height,
		RandaoPartialSig:        NewPartialSigContainer(quorum),
		PostConsensusPartialSig: NewPartialSigContainer(quorum),
		Finished:                false,
	}
}

// ReconstructRandaoSig aggregates collected partial randao sigs, reconstructs a valid sig and returns it
func (pcs *DutyExecutionState) ReconstructRandaoSig(validatorPubKey []byte) ([]byte, error) {
	panic("implement")
}

// ReconstructAttestationSig aggregates collected partial sigs, reconstructs a valid sig and returns an attestation obj with the reconstructed sig
func (pcs *DutyExecutionState) ReconstructAttestationSig(validatorPubKey []byte) (*spec.Attestation, error) {
	// Reconstruct signatures
	signature, err := pcs.PostConsensusPartialSig.ReconstructSignature(validatorPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not reconstruct attestation sig")
	}

	blsSig := spec.BLSSignature{}
	copy(blsSig[:], signature)
	pcs.SignedAttestation.Signature = blsSig
	return pcs.SignedAttestation, nil
}

// SetFinished will mark this execution state as finished
func (pcs *DutyExecutionState) SetFinished() {
	pcs.Finished = true
}

// GetRoot returns the root used for signing and verification
func (pcs *DutyExecutionState) GetRoot() ([]byte, error) {
	marshaledRoot, err := pcs.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode DutyExecutionState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// IsFinished returns true if this execution state is finished
func (pcs *DutyExecutionState) IsFinished() bool {
	return pcs.Finished
}

// Encode returns the encoded struct in bytes or error
func (pcs *DutyExecutionState) Encode() ([]byte, error) {
	return json.Marshal(pcs)
}

// Decode returns error if decoding failed
func (pcs *DutyExecutionState) Decode(data []byte) error {
	return json.Unmarshal(data, &pcs)
}
