package testingutils

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/duty"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var BaseRunner = func() *duty.Runner {
	ret := duty.NewDutyRunner(
		beacon.RoleTypeAttester,
		TestingShare,
		NewTestingQBFTController([]byte{1, 2, 3, 4}),
		NewTestingStorage(),
	)
	ret.StartNewDuty(TestingAttesterDuty)
	return ret
}

var DecidedRunner = func() *duty.Runner {
	return decideRunner(TestAttesterConsensusDataByts, qbft.FirstHeight)
}

var DecidedRunnerWithHeight = func(height qbft.Height) *duty.Runner {
	return decideRunner(TestAttesterConsensusDataByts, height)
}

var DecidedRunnerUnknownDutyType = func() *duty.Runner {
	return decideRunner(TestConsensusUnkownDutyTypeDataByts, qbft.FirstHeight)
}

var decideRunner = func(consensusData []byte, height qbft.Height) *duty.Runner {
	v := BaseValidator()
	for h := qbft.Height(qbft.FirstHeight); h <= height; h++ {
		msgs := []*types.SSVMessage{
			SSVMsg(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       ProposalDataBytes(consensusData, nil, nil),
			}), nil),
			SSVMsg(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       PrepareDataBytes(consensusData),
			}), nil),
			SSVMsg(SignQBFTMsg(TestingSK2, 2, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       PrepareDataBytes(consensusData),
			}), nil),
			SSVMsg(SignQBFTMsg(TestingSK3, 3, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       PrepareDataBytes(consensusData),
			}), nil),
			SSVMsg(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       CommitDataBytes(consensusData),
			}), nil),
			SSVMsg(SignQBFTMsg(TestingSK2, 2, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       CommitDataBytes(consensusData),
			}), nil),
			SSVMsg(SignQBFTMsg(TestingSK3, 3, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       CommitDataBytes(consensusData),
			}), nil),
		}

		if err := v.DutyRunners[beacon.RoleTypeAttester].Decide(TestAttesterConsensusDataByts); err != nil {
			panic(err.Error())
		}
		for _, msg := range msgs {
			if err := v.ProcessMessage(msg); err != nil {
				panic(err.Error())
			}
		}
	}

	return v.DutyRunners[beacon.RoleTypeAttester]
}
