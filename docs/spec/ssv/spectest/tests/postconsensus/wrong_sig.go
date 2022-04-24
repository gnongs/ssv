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
		PostDutyRunnerStateRoot: "33fd61d17dc89513774a7b566e9dddad28ea5703f83efae63aea69e369c1f367",
		ExpectedError:           "partial sig invalid: failed to verify PartialSignature: failed to verify signature",
	}
}
