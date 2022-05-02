package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// ProcessMessage processes network Message of all types
func (v *Validator) ProcessMessage(msg *types.SSVMessage) error {
	dutyRunner := v.DutyRunners.DutyRunnerForMsgID(msg.GetID())
	if dutyRunner == nil {
		return errors.Errorf("could not get duty runner for msg ID")
	}

	if err := v.validateMessage(dutyRunner, msg); err != nil {
		return errors.Wrap(err, "Message invalid")
	}

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		signedMsg := &qbft.SignedMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}
		return v.processConsensusMsg(dutyRunner, signedMsg)
	case types.SSVDecidedMsgType:
		decidedMsg := &qbft.DecidedMessage{}
		if err := decidedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get decided Message from network Message")
		}
		return v.processConsensusMsg(dutyRunner, decidedMsg.SignedMessage)
	case types.SSVPartialSignatureMsgType:
		signedMsg := &SignedPartialSignatureMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		if signedMsg.Message.Type == RandaoPartialSig {
			return v.processRandaoPartialSig(dutyRunner, signedMsg)
		}
		if signedMsg.Message.Type == SelectionProofPartialSig {
			return v.processSelectionProofPartialSig(dutyRunner, signedMsg)
		}
		return v.processPostConsensusSig(dutyRunner, signedMsg)
	default:
		return errors.New("unknown msg")
	}
}

func (v *Validator) validateMessage(runner *Runner, msg *types.SSVMessage) error {
	if runner.CurrentDuty == nil {
		return errors.New("no running duty")
	}

	if !v.share.ValidatorPubKey.MessageIDBelongs(msg.GetID()) {
		return errors.New("msg ID doesn't match validator ID")
	}

	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}

func (v *Validator) processConsensusMsg(dutyRunner *Runner, msg *qbft.SignedMessage) error {
	decided, decidedValue, err := dutyRunner.ProcessConsensusMessage(msg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	postConsensusMsg, err := dutyRunner.SignDutyPostConsensus(decidedValue, v.signer)
	if err != nil {
		return errors.Wrap(err, "failed to decide duty at runner")
	}

	signedMsg, err := v.signPostConsensusMsg(postConsensusMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign post consensus msg")
	}

	data, err := signedMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode post consensus signature msg")
	}

	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(v.share.ValidatorPubKey, dutyRunner.BeaconRoleType),
		Data:    data,
	}

	if err := v.network.Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}

func (v *Validator) processPostConsensusSig(dutyRunner *Runner, signedMsg *SignedPartialSignatureMessage) error {
	quorum, err := dutyRunner.ProcessPostConsensusMessage(signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	switch dutyRunner.BeaconRoleType {
	case types.BNRoleAttester:
		att, err := dutyRunner.State.ReconstructAttestationSig(v.share.ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus sig")
		}
		if err := v.beacon.SubmitAttestation(att); err != nil {
			return errors.Wrap(err, "could not submit to beacon chain reconstructed attestation")
		}
	case types.BNRoleProposer:
		blk, err := dutyRunner.State.ReconstructBeaconBlockSig(v.share.ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus sig")
		}
		if err := v.beacon.SubmitBeaconBlock(blk); err != nil {
			return errors.Wrap(err, "could not submit to beacon chain reconstructed signed beacon block")
		}
	case types.BNRoleAggregator:
		msg, err := dutyRunner.State.ReconstructSignedAggregateSelectionProofSig(v.share.ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus sig")
		}
		if err := v.beacon.SubmitSignedAggregateSelectionProof(msg); err != nil {
			return errors.Wrap(err, "could not submit to beacon chain reconstructed signed aggregate")
		}
	case types.BNRoleSyncCommittee:
		msg, err := dutyRunner.State.ReconstructSyncCommitteeSig(v.share.ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus sig")
		}
		if err := v.beacon.SubmitSyncMessage(msg); err != nil {
			return errors.Wrap(err, "could not submit to beacon chain reconstructed signed sync committee")
		}
	default:
		return errors.Errorf("unknown duty post consensus sig %s", dutyRunner.BeaconRoleType.String())
	}
	return nil
}

func (v *Validator) processRandaoPartialSig(dutyRunner *Runner, signedMsg *SignedPartialSignatureMessage) error {
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

	if err := dutyRunner.Decide(input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (v *Validator) processSelectionProofPartialSig(dutyRunner *Runner, signedMsg *SignedPartialSignatureMessage) error {
	quorum, err := dutyRunner.ProcessSelectionProofMessage(signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing selection proof message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// reconstruct selection proof sig
	fullSig, err := dutyRunner.State.ReconstructSelectionProofSig(v.share.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "could not reconstruct selection proof sig")
	}

	duty := dutyRunner.CurrentDuty

	// TODO waitToSlotTwoThirds

	// get block data
	res, err := v.beacon.SubmitAggregateSelectionProof(duty.Slot, duty.CommitteeIndex, fullSig)
	if err != nil {
		return errors.Wrap(err, "failed to submit aggregate and proof")
	}

	input := &types.ConsensusData{
		Duty:              duty,
		AggregateAndProof: res,
	}

	if err := dutyRunner.Decide(input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (v *Validator) signPostConsensusMsg(msg *PartialSignatureMessage) (*SignedPartialSignatureMessage, error) {
	signature, err := v.signer.SignRoot(msg, types.PartialSignatureType, v.share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign PartialSignatureMessage for PostConsensusPartialSig")
	}

	return &SignedPartialSignatureMessage{
		Message:   msg,
		Signature: signature,
		Signers:   []types.OperatorID{v.share.OperatorID},
	}, nil
}
