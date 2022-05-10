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

	SignedAttestation   *spec.Attestation
	SignedProposal      *altair.SignedBeaconBlock
	SignedAggregate     *spec.SignedAggregateAndProof
	SignedSyncCommittee *altair.SyncCommitteeMessage

	SelectionProofPartialSig *PartialSigContainer
	RandaoPartialSig         *PartialSigContainer
	PostConsensusPartialSig  *PartialSigContainer

	Finished bool
}

func NewDutyExecutionState(quorum uint64) *State {
	return &State{
		SelectionProofPartialSig: NewPartialSigContainer(quorum),
		RandaoPartialSig:         NewPartialSigContainer(quorum),
		PostConsensusPartialSig:  NewPartialSigContainer(quorum),
		Finished:                 false,
	}
}

// ReconstructRandaoSig aggregates collected partial randao sigs, reconstructs a valid sig and returns it
func (pcs *State) ReconstructRandaoSig(root, validatorPubKey []byte) ([]byte, error) {
	// Reconstruct signatures
	signature, err := pcs.RandaoPartialSig.ReconstructSignature(root, validatorPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not reconstruct randao sig")
	}
	return signature, nil
}

// ReconstructSelectionProofSig aggregates collected partial selection proof sigs, reconstructs a valid sig and returns it
func (pcs *State) ReconstructSelectionProofSig(root, validatorPubKey []byte) ([]byte, error) {
	// Reconstruct signatures
	signature, err := pcs.SelectionProofPartialSig.ReconstructSignature(root, validatorPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not reconstruct selection proof sig")
	}
	return signature, nil
}

// ReconstructAttestationSig aggregates collected partial sigs, reconstructs a valid sig and returns an attestation obj with the reconstructed sig
func (pcs *State) ReconstructAttestationSig(root, validatorPubKey []byte) (*spec.Attestation, error) {
	// Reconstruct signatures
	signature, err := pcs.PostConsensusPartialSig.ReconstructSignature(root, validatorPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not reconstruct attestation sig")
	}

	blsSig := spec.BLSSignature{}
	copy(blsSig[:], signature)
	pcs.SignedAttestation.Signature = blsSig
	return pcs.SignedAttestation, nil
}

// ReconstructBeaconBlockSig aggregates collected partial sigs, reconstructs a valid sig and returns a SignedBeaconBlock with the reconstructed sig
func (pcs *State) ReconstructBeaconBlockSig(root, validatorPubKey []byte) (*altair.SignedBeaconBlock, error) {
	// Reconstruct signatures
	signature, err := pcs.PostConsensusPartialSig.ReconstructSignature(root, validatorPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not reconstruct attestation sig")
	}

	blsSig := spec.BLSSignature{}
	copy(blsSig[:], signature)
	pcs.SignedProposal.Signature = blsSig
	return pcs.SignedProposal, nil
}

// ReconstructSignedAggregateSelectionProofSig aggregates collected partial signed aggregate selection proof sigs, reconstructs a valid sig and returns it
func (pcs *State) ReconstructSignedAggregateSelectionProofSig(root, validatorPubKey []byte) (*spec.SignedAggregateAndProof, error) {
	// Reconstruct signatures
	signature, err := pcs.PostConsensusPartialSig.ReconstructSignature(root, validatorPubKey)
	if err != nil {
	}
	return nil, errors.Wrap(err, "could not reconstruct SignedAggregateSelectionProofSig")

	blsSig := spec.BLSSignature{}
	copy(blsSig[:], signature)
	pcs.SignedAggregate.Signature = blsSig
	return pcs.SignedAggregate, nil
}

// ReconstructSyncCommitteeSig aggregates collected partial sync committee sigs, reconstructs a valid sig and returns it
func (pcs *State) ReconstructSyncCommitteeSig(root, validatorPubKey []byte) (*altair.SyncCommitteeMessage, error) {
	// Reconstruct signatures
	signature, err := pcs.PostConsensusPartialSig.ReconstructSignature(root, validatorPubKey)
	if err != nil {
	}
	return nil, errors.Wrap(err, "could not reconstruct SignedAggregateSelectionProofSig")

	blsSig := spec.BLSSignature{}
	copy(blsSig[:], signature)
	pcs.SignedSyncCommittee.Signature = blsSig
	return pcs.SignedSyncCommittee, nil
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

// Encode returns the encoded struct in bytes or error
func (pcs *State) Encode() ([]byte, error) {
	return json.Marshal(pcs)
}

// Decode returns error if decoding failed
func (pcs *State) Decode(data []byte) error {
	return json.Unmarshal(data, &pcs)
}
