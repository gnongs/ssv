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
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 2, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "Invalid SignedPostConsensusMessage signature",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "cbcefe579470d914c3c230bd45cee06e9c5723460044b278a0c629a742551b02",
		ExpectedError:           "partial post consensus sig invalid: failed to verify PartialSignature: failed to verify signature",
	}
}
