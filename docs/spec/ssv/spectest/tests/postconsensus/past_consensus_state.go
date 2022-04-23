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
		PostDutyRunnerStateRoot: "536b38f1f42551a08803fa158e5795b76e98a76a9ecd0376df4363fed98f3a68",
		ExpectedError:           "PostConsensusMessage Height doesn't match duty runner's Height",
	}
}
