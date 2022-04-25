package attestations

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// DutySlotNotMatchingAttestationSlot tests that a duty slot = attestation slot
func DutySlotNotMatchingAttestationSlot() *tests.SpecTest {
	dr := testingutils.BaseRunner()

	consensusData := &types.ConsensusData{
		Duty: &beacon.Duty{
			Type:                    beacon.RoleTypeAttester,
			PubKey:                  testingutils.TestingValidatorPubKey,
			Slot:                    13,
			ValidatorIndex:          1,
			CommitteeIndex:          3,
			CommitteesAtSlot:        36,
			CommitteeLength:         128,
			ValidatorCommitteeIndex: 11,
		},
		AttestationData: testingutils.TestingAttestationData,
	}
	startingValue, _ := consensusData.Encode()

	// the starting value is not the same as the actual proposal!
	dr.NewExecutionState()
	if err := dr.StartNewConsensusInstance(testingutils.TestAttesterConsensusDataByts); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes(startingValue, nil, nil),
		}), nil),
	}

	return &tests.SpecTest{
		Name:                    "duty slot matches attestation slot",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "039e927a1858548cc411afe6442ee7222d285661058621bdb1e02405c0f344d4",
		ExpectedError:           "failed to process consensus msg: could not process msg: proposal invalid: proposal not justified: proposal value invalid: attestation data slot != duty slot",
	}
}
