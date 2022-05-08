package synccommittee

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// SevenOperators tests a full valcheck + post valcheck + duty sig reconstruction flow for 7 operators
func SevenOperators() *tests.SpecTest {
	ks := testingutils.Testing7SharesSet()
	dr := testingutils.SyncCommitteeRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeConsensusDataByts, nil, nil),
		}), nil),

		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),

		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommittee(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
		}), nil),

		testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
		testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[4], 4)),
		testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[5], 5)),
	}

	return &tests.SpecTest{
		Name:                    "sync committee 7 operators happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestingSyncCommitteeDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "82d2d1787358f78e30043586166d906479da3238c6258ef8d962b4601b16e8c0",
	}
}
