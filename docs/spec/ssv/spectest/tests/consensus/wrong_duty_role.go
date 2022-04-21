package consensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongDutyRole tests decided value with duty role != the duty runner's role
func WrongDutyRole() *tests.SpecTest {
	dr := testingutils.BaseRunner()
	if err := dr.StartNewInstance([]byte{1, 2, 3, 4}); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes(testingutils.TestAggregatorConsensusDataByts, nil, nil),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
	}

	return &tests.SpecTest{
		Name:                    "wrong decided value's duty type",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "ffa08170540cdb023aab2321a74d2af8779cad36e2ed0b8e29ee3160b975cf24",
		ExpectedError:           "decided value is invalid: decided value's duty has wrong beacon role type",
	}
}
