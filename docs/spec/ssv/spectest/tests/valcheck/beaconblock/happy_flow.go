package beaconblock

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// HappyFlow tests a full valcheck + post valcheck + duty sig reconstruction flow
func HappyFlow() *tests.SpecTest {
	dr := testingutils.ProposerRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(testingutils.TestingSK1, 1)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(testingutils.TestingSK2, 2)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(testingutils.TestingSK3, 3)),

		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestProposerConsensusDataByts, nil, nil),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),

		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(testingutils.TestingSK1, 1)),
		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(testingutils.TestingSK2, 2)),
		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(testingutils.TestingSK3, 3)),
	}

	return &tests.SpecTest{
		Name:                    "proposer happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestProposerConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "7be6389c14b1ad08cc9bb923cff20b53b22b1b535795afe17a8df122ef4fa5d4",
	}
}
