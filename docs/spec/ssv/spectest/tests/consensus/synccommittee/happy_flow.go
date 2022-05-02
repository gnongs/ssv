package synccommittee

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// HappyFlow tests a full valcheck + post valcheck + duty sig reconstruction flow
func HappyFlow() *tests.SpecTest {
	dr := testingutils.SyncCommitteeRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeConsensusDataByts, nil, nil),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),

		testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(testingutils.TestingSK1, 1)),
		testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(testingutils.TestingSK2, 2)),
		testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(testingutils.TestingSK3, 3)),
	}

	return &tests.SpecTest{
		Name:                    "sync committee happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestingSyncCommitteeDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "b7dfd81271312d188de05aa832c89970efb83f94b7de8d54c6bb5619bfe44a26",
	}
}
