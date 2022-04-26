package duty

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/pkg/errors"
)

// ProcessRandaoMessage process randao msg, returns true if it has quorum for partial signatures.
// returns true only once (first time quorum achieved)
func (dr *Runner) ProcessRandaoMessage(msg *ssv.SignedPartialSignatureMessage) (bool, error) {
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

// canProcessRandaoMsg returns true if it can process randao message, false if not
func (dr *Runner) canProcessRandaoMsg(msg *ssv.SignedPartialSignatureMessage) error {
	if err := dr.validatePartialSigMsg(msg, dr.State.RandaoPartialSig); err != nil {
		return errors.Wrap(err, "randao msg invalid")
	}

	if dr.randaoSigTimeout(dr.BeaconNetwork.EstimatedCurrentSlot()) {
		return errors.New("randao sig collection timeout")
	}

	return nil
}

// randaoSigTimeout returns true if collecting post consensus sigs timed out
func (dr *Runner) randaoSigTimeout(currentSlot spec.Slot) bool {
	return dr.partialSigCollectionTimeout(dr.State.RandaoPartialSig, currentSlot)
}
