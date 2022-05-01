package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (dr *Runner) SignSlotWithSelectionProofPreConsensus(slot spec.Slot, signer types.KeyManager) (*PartialSignatureMessage, error) {
	sig, r, err := signer.SignSlotWithSelectionProof(slot, dr.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign partial selection proof")
	}

	dr.State.SelectionProofPartialSig.SigRoot = ensureRoot(r)

	// generate partial sig for randao
	msg := &PartialSignatureMessage{
		Type:             SelectionProofPartialSig,
		Slot:             slot,
		PartialSignature: sig,
		SigningRoot:      ensureRoot(r),
		Signers:          []types.OperatorID{dr.Share.OperatorID},
	}

	return msg, nil
}

// ProcessSelectionProofMessage process selection proof msg, returns true if it has quorum for partial signatures.
// returns true only once (first time quorum achieved)
func (dr *Runner) ProcessSelectionProofMessage(msg *SignedPartialSignatureMessage) (bool, error) {
	if err := dr.canProcessSelectionProofMsg(msg); err != nil {
		return false, errors.Wrap(err, "can't process selection proof message")
	}

	prevQuorum := dr.State.SelectionProofPartialSig.HasQuorum()

	if err := dr.State.SelectionProofPartialSig.AddSignature(msg.Message); err != nil {
		return false, errors.Wrap(err, "could not add partial selection proof signature")
	}

	if prevQuorum {
		return false, nil
	}

	return dr.State.SelectionProofPartialSig.HasQuorum(), nil
}

// canProcessRandaoMsg returns true if it can process selection proof message, false if not
func (dr *Runner) canProcessSelectionProofMsg(msg *SignedPartialSignatureMessage) error {
	if err := dr.validatePartialSigMsg(msg, dr.State.SelectionProofPartialSig, dr.CurrentDuty.Slot); err != nil {
		return errors.Wrap(err, "selection proof msg invalid")
	}

	if dr.selectionProofSigTimeout(dr.BeaconNetwork.EstimatedCurrentSlot()) {
		return errors.New("selection proof sig collection timeout")
	}

	return nil
}

// selectionProofSigTimeout returns true if collecting selection proof sigs timed out
func (dr *Runner) selectionProofSigTimeout(currentSlot spec.Slot) bool {
	return dr.partialSigCollectionTimeout(dr.State.SelectionProofPartialSig, currentSlot)
}
