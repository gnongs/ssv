package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processConsensusMsg(dutyRunner *DutyRunner, msg *qbft.SignedMessage) error {
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
		MsgID:   types.MessageIDForValidatorPKAndRole(v.share.ValidatorPubKey, dutyRunner.BeaconRoleType),
		Data:    data,
	}

	if err := v.network.Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}
