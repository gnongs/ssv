package ssv

import (
	"bytes"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processConsensusMsg(dutyRunner *DutyRunner, msg *qbft.SignedMessage) error {
	decided, decidedValueByts, err := dutyRunner.QBFTController.ProcessMsg(msg)
	if err != nil {
		return errors.Wrap(err, "failed to process consensus msg")
	}

	/**
	Decided returns true only once so if it is true it must be for the current running instance
	*/
	if !decided {
		return nil
	}

	decidedValue := &types.ConsensusData{}
	if err := decidedValue.Decode(decidedValueByts); err != nil {
		return errors.Wrap(err, "failed to parse decided value to ConsensusData")
	}

	if err := v.validateDecidedConsensusData(dutyRunner, decidedValue); err != nil {
		return errors.Wrap(err, "decided value is invalid")
	}

	postConsensusMsg, err := dutyRunner.DecideRunningInstance(decidedValue, v.signer)
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
		MsgID:   types.MessageIDForValidatorPKAndRole(v.share.ValidatorPubKey, dutyRunner.BeaconRoleType),
		Data:    data,
	}

	if err := v.network.Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial sig")
	}
	return nil
}

func (v *Validator) validateDecidedConsensusData(dutyRunner *DutyRunner, val *types.ConsensusData) error {
	byts, err := val.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode decided value")
	}
	if err := v.valCheck(byts); err != nil {
		return errors.Wrap(err, "decided value is invalid")
	}

	if dutyRunner.BeaconRoleType != val.Duty.Type {
		return errors.New("decided value's duty has wrong beacon role type")
	}

	if !bytes.Equal(dutyRunner.Share.ValidatorPubKey, val.Duty.PubKey[:]) {
		return errors.New("decided value's validator pk is wrong")
	}

	return nil
}
