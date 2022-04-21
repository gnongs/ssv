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
		PostDutyRunnerStateRoot: "0576fbb0e4c7011caa73b31fd0b50bd20a8b997b8ce1a49cb22d25ded3d18642",
		ExpectedError:           "PostConsensusMessage Height doesn't match duty runner's Height",
	}
}
