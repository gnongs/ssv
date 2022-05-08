package attester

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// HappyFlow tests a full valcheck + post valcheck + duty sig reconstruction flow
func HappyFlow() *tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),

		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "attester happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestAttesterConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "4368bab5aed03dfb6a62ee59faf9078912588a892ed21cae551dd51719cc2f4f",
	}
}
