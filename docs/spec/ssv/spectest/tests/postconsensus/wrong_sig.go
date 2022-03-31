package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// InvaliSignature tests an invalid SignedPostConsensusMessage sig
func InvaliSignature() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 2, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "Invalid SignedPostConsensusMessage signature",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "854de580f7c23c94607d671315b57afef2ca8494859ee5f3d4af235ba50c55bb",
		ExpectedError:           "partial sig invalid: failed to verify DutySignature: failed to verify signature",
	}
}
