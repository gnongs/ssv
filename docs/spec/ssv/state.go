package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/altair"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// State holds all the relevant progress the duty execution progress
type State struct {
	RunningInstance *qbft.Instance

	DecidedValue *types.ConsensusData

	SignedAttestation *spec.Attestation
	SignedProposal    *altair.SignedBeaconBlock

	RandaoPartialSig        *PartialSigContainer
	PostConsensusPartialSig *PartialSigContainer

	Finished bool
}

func NewDutyExecutionState(quorum uint64) *State {
	return &State{
		RandaoPartialSig:        NewPartialSigContainer(quorum),
		PostConsensusPartialSig: NewPartialSigContainer(quorum),
		Finished:                false,
	}
}

// ReconstructRandaoSig aggregates collected partial randao sigs, reconstructs a valid sig and returns it
func (pcs *State) ReconstructRandaoSig(validatorPubKey []byte) ([]byte, error) {
	// Reconstruct signatures
	signature, err := pcs.RandaoPartialSig.ReconstructSignature(validatorPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not reconstruct randao sig")
	}
	return signature, nil
}

// ReconstructAttestationSig aggregates collected partial sigs, reconstructs a valid sig and returns an attestation obj with the reconstructed sig
func (pcs *State) ReconstructAttestationSig(validatorPubKey []byte) (*spec.Attestation, error) {
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

// ReconstructBeaconBlockSig aggregates collected partial sigs, reconstructs a valid sig and returns a SignedBeaconBlock with the reconstructed sig
func (pcs *State) ReconstructBeaconBlockSig(validatorPubKey []byte) (*altair.SignedBeaconBlock, error) {
	// Reconstruct signatures
	signature, err := pcs.PostConsensusPartialSig.ReconstructSignature(validatorPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not reconstruct attestation sig")
	}

	blsSig := spec.BLSSignature{}
	copy(blsSig[:], signature)
	pcs.SignedProposal.Signature = blsSig
	return pcs.SignedProposal, nil
}

// SetFinished will mark this execution state as finished
func (pcs *State) SetFinished() {
	pcs.Finished = true
}

// GetRoot returns the root used for signing and verification
func (pcs *State) GetRoot() ([]byte, error) {
	marshaledRoot, err := pcs.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode State")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// IsFinished returns true if this execution state is finished
func (pcs *State) IsFinished() bool {
	return pcs.Finished
}

// Encode returns the encoded struct in bytes or error
func (pcs *State) Encode() ([]byte, error) {
	return json.Marshal(pcs)
}

// Decode returns error if decoding failed
func (pcs *State) Decode(data []byte) error {
	return json.Unmarshal(data, &pcs)
}
