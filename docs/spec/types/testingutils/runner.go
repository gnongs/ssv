package testingutils

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var AttesterRunner = func() *ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.BeaconAttestationValueCheck(NewTestingKeyManager(), ssv.NowTestNetwork), 4)
}

var AttesterRunner7Operators = func() *ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.BeaconAttestationValueCheck(NewTestingKeyManager(), ssv.NowTestNetwork), 7)
}

var ProposerRunner = func() *ssv.Runner {
	return baseRunner(types.BNRoleProposer, ssv.BeaconBlockValueCheck(NewTestingKeyManager(), ssv.NowTestNetwork), 4)
}

var AggregatorRunner = func() *ssv.Runner {
	return baseRunner(types.BNRoleAggregator, ssv.AggregatorValueCheck(NewTestingKeyManager(), ssv.NowTestNetwork), 4)
}

var SyncCommitteeRunner = func() *ssv.Runner {
	return baseRunner(types.BNRoleSyncCommittee, ssv.SyncCommitteeValueCheck(NewTestingKeyManager(), ssv.NowTestNetwork), 4)
}

var baseRunner = func(role types.BeaconRole, valCheck qbft.ProposedValueCheck, operatorsCnt int) *ssv.Runner {
	share := TestingShare
	if operatorsCnt == 7 {
		share = TestingShareSevenOperators
	}

	return ssv.NewDutyRunner(
		role,
		ssv.NowTestNetwork,
		share,
		NewTestingQBFTController(types.NewMsgID(TestingValidatorPubKey[:], role), valCheck, share),
		NewTestingStorage(),
		valCheck,
	)
}

var DecidedRunner = func() *ssv.Runner {
	return decideRunner(TestAttesterConsensusDataByts, qbft.FirstHeight)
}

var DecidedRunnerWithHeight = func(height qbft.Height) *ssv.Runner {
	return decideRunner(TestAttesterConsensusDataByts, height)
}

var DecidedRunnerUnknownDutyType = func() *ssv.Runner {
	return decideRunner(TestConsensusUnkownDutyTypeDataByts, qbft.FirstHeight)
}

var decideRunner = func(consensusData []byte, height qbft.Height) *ssv.Runner {
	v := BaseValidator()
	for h := qbft.Height(qbft.FirstHeight); h <= height; h++ {
		msgs := []*types.SSVMessage{
			SSVMsgAttester(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       ProposalDataBytes(consensusData, nil, nil),
			}), nil),
			SSVMsgAttester(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       PrepareDataBytes(consensusData),
			}), nil),
			SSVMsgAttester(SignQBFTMsg(TestingSK2, 2, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       PrepareDataBytes(consensusData),
			}), nil),
			SSVMsgAttester(SignQBFTMsg(TestingSK3, 3, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       PrepareDataBytes(consensusData),
			}), nil),
			SSVMsgAttester(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       CommitDataBytes(consensusData),
			}), nil),
			SSVMsgAttester(SignQBFTMsg(TestingSK2, 2, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       CommitDataBytes(consensusData),
			}), nil),
			SSVMsgAttester(SignQBFTMsg(TestingSK3, 3, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       CommitDataBytes(consensusData),
			}), nil),
		}

		if err := v.DutyRunners[types.BNRoleAttester].Decide(TestAttesterConsensusData); err != nil {
			panic(err.Error())
		}
		for _, msg := range msgs {
			if err := v.ProcessMessage(msg); err != nil {
				panic(err.Error())
			}
		}
	}

	return v.DutyRunners[types.BNRoleAttester]
}
