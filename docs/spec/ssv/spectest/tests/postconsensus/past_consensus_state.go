package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// PastConsensusState tests msg for a past consensus state
func PastConsensusState() *tests.SpecTest {
	dr := testingutils.DecidedRunnerWithHeight(1)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "past SignedPostConsensusMessage",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "3547b8872a38ec369386d127767979ccb48e7ce28d5a4db22502c97bf768994f",
		ExpectedError:           "PostConsensusMessage Height doesn't match duty runner's Height",
	}
}
