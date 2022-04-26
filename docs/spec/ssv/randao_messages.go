package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processRandaoPartialSig(dutyRunner *DutyRunner, signedMsg *SignedPartialSignatureMessage) error {
	quorum, err := dutyRunner.ProcessRandaoMessage(signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing randao message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// randao is relevant only for block proposals, no need to check type
	fullSig, err := dutyRunner.State.ReconstructRandaoSig(v.share.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "could not reconstruct randao sig")
	}

	duty := dutyRunner.CurrentDuty

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
