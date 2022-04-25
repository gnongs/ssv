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
		PostDutyRunnerStateRoot: "0351bb303531bc5858d312928c9577e9ca0104f3d8986a34fce30f2519908b1e",
	}
}
