package aggregator

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// HappyFlow tests a full valcheck + post valcheck + duty sig reconstruction flow
func HappyFlow() *tests.SpecTest {
	dr := testingutils.AggregatorRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(testingutils.TestingSK1, 1)),
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(testingutils.TestingSK2, 2)),
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(testingutils.TestingSK3, 3)),

		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestAggregatorConsensusDataByts, nil, nil),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),

		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(testingutils.TestingSK1, 1)),
		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(testingutils.TestingSK2, 2)),
		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(testingutils.TestingSK3, 3)),
	}

	return &tests.SpecTest{
		Name:                    "aggregator happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestAggregatorConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "8ff71c2e43b1e8df40af0d7412716690dfc73ecc785763d1d77ebd7239fd237f",
	}
}
