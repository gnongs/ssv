package consensus

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

var testConsensusAttDataNil = &types.ConsensusData{
	Duty: &beacon.Duty{
		Type:                    beacon.RoleTypeAttester,
		PubKey:                  testingutils.TestingValidatorPubKey,
		Slot:                    12,
		ValidatorIndex:          1,
		CommitteeIndex:          22,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	},
	AttestationData: nil,
}
var testConsensusAttDataNilByts, _ = testConsensusAttDataNil.Encode()

// AttestationDataNil tests decided value attestation data nil (if duty beacon role is attester)
func AttestationDataNil() *tests.SpecTest {
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
			Data:       testingutils.ProposalDataBytes(testConsensusAttDataNilByts, nil, nil),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testConsensusAttDataNilByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testConsensusAttDataNilByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testConsensusAttDataNilByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testConsensusAttDataNilByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testConsensusAttDataNilByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testConsensusAttDataNilByts),
		}), nil),
	}

	return &tests.SpecTest{
		Name:                    "decided value's attestation data nil (role attester)",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "ed2563907d5359a979be22f89d9f8d858394df17bde3a8b24c9c2b848ffb6d98",
		ExpectedError:           "decided value is invalid: decided value's AttestationData is nil",
	}
}
