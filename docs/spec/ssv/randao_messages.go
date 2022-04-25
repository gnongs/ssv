package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processRandaoPartialSig(dutyRunner *DutyRunner, signedMsg *SignedPartialSignatureMessage) error {
	postCons := dutyRunner.PostConsensusStateForHeight(signedMsg.Message.Height)
	if postCons == nil {
		return errors.New("PartialSignatureMessage Height doesn't match duty runner's Height")
	}

	if err := v.validateRandaoPartialSig(postCons, signedMsg); err != nil {
		return errors.Wrap(err, "partial randao sig invalid")
	}

	prevQuorum := postCons.RandaoPartialSig.HasQuorum()

	if err := postCons.RandaoPartialSig.AddSignature(signedMsg.Message); err != nil {
		return errors.Wrap(err, "could not add partial randao signature")
	}

	if prevQuorum || !postCons.RandaoPartialSig.HasQuorum() {
		return nil
	}

	// randao is relevant only for block proposals, no need to check type
	fullSig, err := postCons.ReconstructRandaoSig(v.share.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "could not reconstruct randao sig")
	}

	duty := dutyRunner.DutyExecutionState.ProposedValue.Duty

	// get block data
	blk, err := v.beacon.GetBeaconBlock(duty.Slot, duty.CommitteeIndex, v.share.Graffiti, fullSig)
	if err != nil {
		return errors.Wrap(err, "failed to get beacon block")
	}

	input := &types.ConsensusData{
		Duty:      duty,
		BlockData: blk,
	}

	byts, err := input.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode input")
	}

	// validate input
	//we should maybe move the val check to the duty runner as it needs to change with each duty type
	if err := v.valCheck(byts); err != nil {
		return errors.Wrap(err, "StartDuty input data invalid")
	}

	if err := dutyRunner.StartNewConsensusInstance(byts); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (v *Validator) validateRandaoPartialSig(executionState *DutyExecutionState, SignedMsg *SignedPartialSignatureMessage) error {
	panic("implement")
}
