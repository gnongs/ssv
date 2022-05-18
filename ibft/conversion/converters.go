package conversion

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/network"
	"github.com/bloxapp/ssv/protocol/v1/message"
	"github.com/bloxapp/ssv/utils/format"
	"github.com/bloxapp/ssv/utils/logex"
)

// converters are used to encapsulate the struct of the messages
// that are passed in the network (v0), until v1 fork

// ToV1Message converts an old message to v1
func ToV1Message(msgV0 *network.Message) (*message.SSVMessage, error) {
	msg := message.SSVMessage{}

	switch msgV0.Type {
	case network.NetworkMsg_SyncType:
		msg.MsgType = message.SSVSyncMsgType
		if msgV0.SyncMessage != nil {
			identifier := toIdentifierV1(msgV0.SyncMessage.GetLambda())
			msg.ID = identifier
			syncMsg := new(message.SyncMessage)
			syncMsg.Status = message.StatusSuccess
			syncMsg.Params = new(message.SyncParams)
			syncMsg.Params.Identifier = identifier
			params := msgV0.SyncMessage.GetParams()
			syncMsg.Params.Height = make([]message.Height, 0)
			for _, p := range params {
				syncMsg.Params.Height = append(syncMsg.Params.Height, message.Height(p))
			}
			for _, sm := range msgV0.SyncMessage.GetSignedMessages() {
				signed, err := ToSignedMessageV1(sm)
				if err != nil {
					return nil, err
				}
				if signed.Message != nil {
					syncMsg.Data = append(syncMsg.Data, signed)
				}
			}
			if len(syncMsg.Data) == 0 {
				syncMsg.Status = message.StatusNotFound
			}
			if err := msgV0.SyncMessage.Error; len(err) > 0 {
				logex.GetLogger().Warn("sync message error", zap.String("err", err))
				syncMsg.Status = message.StatusError
			}
			switch msgV0.SyncMessage.Type {
			case network.Sync_GetHighestType:
				syncMsg.Protocol = message.LastDecidedType
			case network.Sync_GetLatestChangeRound:
				syncMsg.Protocol = message.LastChangeRoundType
			case network.Sync_GetInstanceRange:
				syncMsg.Protocol = message.DecidedHistoryType
			}
			data, err := json.Marshal(syncMsg)
			if err != nil {
				return nil, err
			}
			msg.Data = data
			//msg.ID = msgV0.SyncMessage.GetLambda()
		}
		return &msg, nil
	case network.NetworkMsg_SignatureType:
		msg.MsgType = message.SSVPostConsensusMsgType
		postConMsg := toSignedPostConsensusMessageV1(msgV0.SignedMessage)
		data, err := postConMsg.Encode()
		if err != nil {
			return nil, err
		}
		msg.Data = data
		if identifierV0 := msgV0.SignedMessage.GetMessage().GetLambda(); identifierV0 != nil {
			msg.ID = toIdentifierV1(identifierV0)
		}
		return &msg, nil
	case network.NetworkMsg_IBFTType:
		msg.MsgType = message.SSVConsensusMsgType
	case network.NetworkMsg_DecidedType:
		msg.MsgType = message.SSVDecidedMsgType
	}

	if msgV0.SignedMessage != nil {
		signed, err := ToSignedMessageV1(msgV0.SignedMessage)
		if err != nil {
			return nil, err
		}
		data, err := signed.Encode()
		if err != nil {
			return nil, err
		}
		msg.Data = data
		if len(msg.ID) == 0 {
			msg.ID = signed.Message.Identifier
		}
	}

	return &msg, nil
}

func toSignedPostConsensusMessageV1(sm *proto.SignedMessage) *message.SignedPostConsensusMessage {
	signed := new(message.SignedPostConsensusMessage)
	consensus := &message.PostConsensusMessage{
		Height:          message.Height(sm.Message.SeqNumber),
		DutySignature:   sm.GetSignature(),
		DutySigningRoot: nil,
	}

	var signers []message.OperatorID
	for _, signer := range sm.GetSignerIds() {
		signers = append(signers, message.OperatorID(signer))
	}
	consensus.Signers = signers

	signed.Message = consensus
	signed.Signers = signers
	signed.Signature = sm.GetSignature() // TODO should be message sign and not duty sign

	return signed
}

// ToSignedMessageV1 converts a signed message from v0 to v1
func ToSignedMessageV1(sm *proto.SignedMessage) (*message.SignedMessage, error) {
	signed := new(message.SignedMessage)
	signed.Signature = sm.GetSignature()
	signers := sm.GetSignerIds()
	for _, s := range signers {
		signed.Signers = append(signed.Signers, message.OperatorID(s))
	}
	if msg := sm.GetMessage(); msg != nil {
		signed.Message = new(message.ConsensusMessage)
		data := msg.GetValue()
		signed.Message.Round = message.Round(msg.GetRound())
		signed.Message.Identifier = toIdentifierV1(msg.GetLambda())
		signed.Message.Height = message.Height(msg.GetSeqNumber())
		switch msg.GetType() {
		case proto.RoundState_NotStarted:
			// TODO
		case proto.RoundState_PrePrepare:
			signed.Message.MsgType = message.ProposalMsgType
			p, err := (&message.ProposalData{Data: data}).Encode()
			if err != nil {
				return nil, err
			}
			signed.Message.Data = p
		case proto.RoundState_Prepare:
			signed.Message.MsgType = message.PrepareMsgType
			p, err := (&message.PrepareData{Data: data}).Encode()
			if err != nil {
				return nil, err
			}
			signed.Message.Data = p
		case proto.RoundState_Commit:
			signed.Message.MsgType = message.CommitMsgType
			c, err := (&message.CommitData{Data: data}).Encode()
			if err != nil {
				return nil, err
			}
			signed.Message.Data = c
		case proto.RoundState_ChangeRound:
			signed.Message.MsgType = message.RoundChangeMsgType
			rcd, err := toV1ChangeRound(data)
			if err != nil {
				return nil, err
			}
			signed.Message.Data = rcd
		case proto.RoundState_Stopped:
			// TODO
		}
	}
	return signed, nil
}

func toV1ChangeRound(changeRoundData []byte) ([]byte, error) {
	// TODO need to remove log once done with testing
	r, err := json.Marshal(changeRoundData)
	if err == nil {
		logex.GetLogger().Debug("------ convert change round v0 -> v1", zap.String("data marshaled", string(r)), zap.ByteString("data byte", changeRoundData))
	} else {
		logex.GetLogger().Debug("------ FAILED convert change round v0 -> v1", zap.Error(err))
	}
	ret := &proto.ChangeRoundData{}
	if err := json.Unmarshal(changeRoundData, ret); err != nil {
		logex.GetLogger().Warn("failed to unmarshal v0 change round struct", zap.Error(err))
		return new(message.RoundChangeData).Encode() // should return empty struct
	}

	var signers []message.OperatorID
	for _, signer := range ret.GetSignerIds() {
		signers = append(signers, message.OperatorID(signer))
	}

	consensusMsg := &message.ConsensusMessage{}
	if ret.GetJustificationMsg() != nil {
		consensusMsg.Height = message.Height(ret.GetJustificationMsg().SeqNumber)
		consensusMsg.Round = message.Round(ret.GetJustificationMsg().Round)
		consensusMsg.Identifier = toIdentifierV1(ret.GetJustificationMsg().Lambda)
		consensusMsg.Data = ret.GetJustificationMsg().Value
		consensusMsg.MsgType = message.PrepareMsgType // can be only prepare
	}

	crm := &message.RoundChangeData{
		PreparedValue:    ret.GetPreparedValue(),
		Round:            message.Round(ret.GetPreparedRound()),
		NextProposalData: nil,
		RoundChangeJustification: []*message.SignedMessage{{
			Signature: ret.GetJustificationSig(),
			Signers:   signers,
			Message:   consensusMsg,
		}},
	}

	encoded, err := crm.Encode()
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// ToV0Message converts v1 message to v0
func ToV0Message(msg *message.SSVMessage) (*network.Message, error) {
	v0Msg := &network.Message{}
	identifierV0 := toIdentifierV0(msg.GetIdentifier())
	if msg.GetType() == message.SSVDecidedMsgType {
		v0Msg.Type = network.NetworkMsg_DecidedType // TODO need to provide the proper type (under consensus or post consensus?)
	}
	switch msg.GetType() {
	case message.SSVConsensusMsgType, message.SSVDecidedMsgType:
		if v0Msg.Type != network.NetworkMsg_DecidedType {
			v0Msg.Type = network.NetworkMsg_IBFTType
		}
		signedMsg := &message.SignedMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return nil, errors.Wrap(err, "could not decode consensus signed message")
		}

		sm, err := ToSignedMessageV0(signedMsg, identifierV0)
		if err != nil {
			return nil, err
		}
		v0Msg.SignedMessage = sm
		switch v0Msg.SignedMessage.GetMessage().GetType() {
		case proto.RoundState_ChangeRound:
			v0Msg.Type = network.NetworkMsg_IBFTType
		case proto.RoundState_Decided:
			v0Msg.Type = network.NetworkMsg_DecidedType
		default:
		}
	case message.SSVPostConsensusMsgType:
		v0Msg.Type = network.NetworkMsg_SignatureType
		signedMsg := &message.SignedPostConsensusMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return nil, errors.Wrap(err, "could not get post consensus Message from network Message")
		}
		v0Msg.SignedMessage = toSignedMessagePostConsensusV0(signedMsg, identifierV0)
	case message.SSVSyncMsgType:
		v0Msg.Type = network.NetworkMsg_SyncType
		syncMsg := &message.SyncMessage{}
		if err := syncMsg.Decode(msg.GetData()); err != nil {
			return nil, errors.Wrap(err, "could not decode consensus signed message")
		}
		if syncMsg.Status == message.StatusUnknown {
			syncMsg.Status = message.StatusSuccess
		}
		v0Msg.SyncMessage = new(network.SyncMessage)
		if len(syncMsg.Params.Height) > 0 {
			v0Msg.SyncMessage.Params = make([]uint64, 1)
			v0Msg.SyncMessage.Params[0] = uint64(syncMsg.Params.Height[0])
			if len(syncMsg.Params.Height) > 1 {
				v0Msg.SyncMessage.Params = append(v0Msg.SyncMessage.Params, uint64(syncMsg.Params.Height[1]))
			}
		}
		v0Msg.SyncMessage.Lambda = identifierV0
		switch syncMsg.Status {
		case message.StatusSuccess:
			v0Msg.SyncMessage.SignedMessages = make([]*proto.SignedMessage, 0)
			for _, smsg := range syncMsg.Data {
				sm, err := ToSignedMessageV0(smsg, identifierV0)
				if err != nil {
					return nil, err
				}
				v0Msg.SyncMessage.SignedMessages = append(v0Msg.SyncMessage.SignedMessages, sm)
			}
		case message.StatusNotFound:
			v0Msg.SyncMessage.SignedMessages = make([]*proto.SignedMessage, 0)
		default:
			v0Msg.SyncMessage.Error = "error"
		}
		switch syncMsg.Protocol {
		case message.LastDecidedType:
			v0Msg.SyncMessage.Type = network.Sync_GetHighestType
		case message.LastChangeRoundType:
			v0Msg.SyncMessage.Type = network.Sync_GetLatestChangeRound
		case message.DecidedHistoryType:
			v0Msg.SyncMessage.Type = network.Sync_GetInstanceRange
		}
	default:
		return nil, errors.New("unknown msg")
	}

	return v0Msg, nil
}

// ToSignedMessageV0 converts a signed message from v1 to v0
func ToSignedMessageV0(signedMsg *message.SignedMessage, identifierV0 []byte) (*proto.SignedMessage, error) {
	signedMsgV0 := &proto.SignedMessage{}
	signedMsgV0.Message = &proto.Message{
		Round:     uint64(signedMsg.Message.Round),
		Lambda:    identifierV0,
		SeqNumber: uint64(signedMsg.Message.Height),
		Value:     nil,
	}

	switch signedMsg.Message.MsgType {
	case message.ProposalMsgType:
		signedMsgV0.Message.Type = proto.RoundState_PrePrepare
		p, err := signedMsg.Message.GetProposalData()
		if err != nil {
			return nil, err
		}
		signedMsgV0.Message.Value = p.Data
	case message.PrepareMsgType:
		signedMsgV0.Message.Type = proto.RoundState_Prepare
		p, err := signedMsg.Message.GetPrepareData()
		if err != nil {
			return nil, err
		}
		signedMsgV0.Message.Value = p.Data
	case message.CommitMsgType:
		signedMsgV0.Message.Type = proto.RoundState_Commit
		c, err := signedMsg.Message.GetCommitData()
		if err != nil {
			return nil, err
		}
		signedMsgV0.Message.Value = c.Data

	case message.RoundChangeMsgType:
		signedMsgV0.Message.Type = proto.RoundState_ChangeRound
		cr, err := signedMsg.Message.GetRoundChangeData()
		if err != nil {
			return nil, err
		}
		if cr.GetPreparedValue() != nil && len(cr.GetPreparedValue()) > 0 {
			signedMsgV0.Message.Value = cr.GetPreparedValue()
		} else {
			v := make(map[string]interface{})
			marshaledV, err := json.Marshal(v)
			if err != nil {
				return nil, errors.Wrap(err, "failed to marshal empty map")
			}
			signedMsgV0.Message.Value = marshaledV // adding empty json in order to support v0 root
		}
	}

	signedMsgV0.Signature = signedMsg.GetSignature()
	for _, signer := range signedMsg.GetSigners() {
		signedMsgV0.SignerIds = append(signedMsgV0.SignerIds, uint64(signer))
	}
	return signedMsgV0, nil
}

func toSignedMessagePostConsensusV0(signedMsg *message.SignedPostConsensusMessage, identifierV0 []byte) *proto.SignedMessage {
	signedMsgV0 := &proto.SignedMessage{}
	signedMsgV0.Message = &proto.Message{
		Lambda:    identifierV0,
		SeqNumber: uint64(signedMsg.Message.Height),
	}
	signedMsgV0.Signature = signedMsg.Message.DutySignature
	for _, signer := range signedMsg.Signers {
		signedMsgV0.SignerIds = append(signedMsgV0.SignerIds, uint64(signer))
	}
	return signedMsgV0
}

func toIdentifierV0(mid message.Identifier) []byte {
	return []byte(format.IdentifierFormat(mid.GetValidatorPK(), mid.GetRoleType().String()))
}

func toIdentifierV1(old []byte) message.Identifier {
	pk, rt := format.IdentifierUnformat(string(old))
	pkraw, err := hex.DecodeString(pk)
	if err != nil {
		return nil
	}
	return message.NewIdentifier(pkraw, message.RoleTypeFromString(rt))
}
