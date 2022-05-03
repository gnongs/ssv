package testingutils

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var TestingConfig = &qbft.Config{
	Signer:     NewTestingKeyManager(),
	SigningPK:  TestingSK1.GetPublicKey().Serialize(),
	Domain:     types.PrimusTestnet,
	ValueCheck: ssv.BeaconAttestationValueCheck(NewTestingKeyManager(), ssv.NowTestNetwork),
	Storage:    NewTestingStorage(),
	Network:    NewTestingNetwork(),
}
var TestingShare = testShare(fourOperatorsCommittee, 3, 2)
var TestingShareSevenOperators = testShare(sevenOperatorsCommittee, 5, 3)

var fourOperatorsCommittee = []*types.Operator{
	{
		OperatorID: 1,
		PubKey:     TestingSK1.GetPublicKey().Serialize(),
	},
	{
		OperatorID: 2,
		PubKey:     TestingSK2.GetPublicKey().Serialize(),
	},
	{
		OperatorID: 3,
		PubKey:     TestingSK3.GetPublicKey().Serialize(),
	},
	{
		OperatorID: 4,
		PubKey:     TestingSK4.GetPublicKey().Serialize(),
	},
}

var sevenOperatorsCommittee = append(fourOperatorsCommittee, []*types.Operator{
	{
		OperatorID: 5,
		PubKey:     TestingSK5.GetPublicKey().Serialize(),
	},
	{
		OperatorID: 6,
		PubKey:     TestingSK6.GetPublicKey().Serialize(),
	},
	{
		OperatorID: 7,
		PubKey:     TestingSK7.GetPublicKey().Serialize(),
	},
}...)

var testShare = func(committee []*types.Operator, quorum, partialQuorum uint64) *types.Share {
	return &types.Share{
		OperatorID:      1,
		ValidatorPubKey: TestingValidatorPubKey[:],
		SharePubKey:     TestingSK1.GetPublicKey().Serialize(),
		DomainType:      types.PrimusTestnet,
		Quorum:          quorum,
		PartialQuorum:   partialQuorum,
		Committee:       committee,
	}
}

var BaseInstance = func() *qbft.Instance {
	return baseInstance(TestingShare, []byte{1, 2, 3, 4})
}

var SevenOperatorsInstance = func() *qbft.Instance {
	return baseInstance(TestingShareSevenOperators, []byte{1, 2, 3, 4})
}

var baseInstance = func(share *types.Share, identifier []byte) *qbft.Instance {
	ret := qbft.NewInstance(TestingConfig, nil, nil)
	ret.State = &qbft.State{
		Share:                           share,
		ID:                              identifier,
		Round:                           qbft.FirstRound,
		Height:                          qbft.FirstHeight,
		LastPreparedRound:               qbft.NoRound,
		LastPreparedValue:               nil,
		ProposalAcceptedForCurrentRound: nil,
	}
	ret.State.ProposeContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.State.PrepareContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.State.CommitContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.State.RoundChangeContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	return ret
}

func NewTestingQBFTController(identifier []byte, valCheck qbft.ProposedValueCheck, share *types.Share) *qbft.Controller {
	ret := qbft.NewController(
		identifier,
		share,
		types.PrimusTestnet,
		NewTestingKeyManager(),
		valCheck,
		NewTestingStorage(),
		NewTestingNetwork(),
	)
	ret.Domain = types.PrimusTestnet
	return ret
}
