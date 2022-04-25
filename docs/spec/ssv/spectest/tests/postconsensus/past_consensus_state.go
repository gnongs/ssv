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
		PostDutyRunnerStateRoot: "b3556d24734587bbd424e9de7dca4efca7bcdb2045ebeb4a44faec0d506fdfef",
		ExpectedError:           "PartialSignatureMessage Height doesn't match duty runner's Height",
	}
}
