package attestations

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// ValidAttestation tests a valid attestation proposal
func ValidAttestation() *tests.SpecTest {
	dr := testingutils.BaseRunner()
	startingValue := testingutils.TestAttesterConsensusDataByts

	dr.NewExecutionState()
	if err := dr.StartNewConsensusInstance(startingValue); err != nil {
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
		Name:                    "valid attestation",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "123d7f93154d9786630b47808b2b636ca2bc6fa476900a8610dc86134c0cc1ac",
	}
}
