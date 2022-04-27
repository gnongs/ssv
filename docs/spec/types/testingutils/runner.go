package testingutils

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var AttesterRunner = func() *ssv.Runner {
	return baseRunner(beacon.RoleTypeAttester, ssv.BeaconAttestationValueCheck(NewTestingKeyManager(), ssv.NowTestNetwork))
}

var ProposerRunner = func() *ssv.Runner {
	return baseRunner(beacon.RoleTypeProposer, ssv.BeaconBlockValueCheck(NewTestingKeyManager(), ssv.NowTestNetwork))
}

var baseRunner = func(role beacon.RoleType, valCheck qbft.ProposedValueCheck) *ssv.Runner {
	return ssv.NewDutyRunner(
		role,
		ssv.NowTestNetwork,
		TestingShare,
		NewTestingQBFTController(types.MessageIDForValidatorPKAndRole(TestingValidatorPubKey[:], role), valCheck),
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

		if err := v.DutyRunners[beacon.RoleTypeAttester].Decide(TestAttesterConsensusData); err != nil {
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
