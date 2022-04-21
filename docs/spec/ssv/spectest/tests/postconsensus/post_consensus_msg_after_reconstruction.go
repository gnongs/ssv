package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// MsgAfterReconstruction tests msg received after partial sig reconstructed and consensus state set to finished
func MsgAfterReconstruction() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 1, qbft.FirstHeight)),
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK2, 2, qbft.FirstHeight)),
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK3, 3, qbft.FirstHeight)),
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK4, 4, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "4th msg after reconstruction",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "d020cbfdeb628ab92857b33a94c7a154118f41dff7e9b8368df54864611f0860",
	}
}
