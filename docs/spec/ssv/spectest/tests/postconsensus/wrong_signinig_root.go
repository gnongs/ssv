package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongSigningRoot tests an invalid PostConsensusMessage SigningRoot
func WrongSigningRoot() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsgWithWrongRoot(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "invalid PostConsensusMessage SigningRoot",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "cbcefe579470d914c3c230bd45cee06e9c5723460044b278a0c629a742551b02",
		ExpectedError:           "partial post consensus sig invalid: SignedPartialSignatureMessage invalid: SigningRoot invalid",
	}
}
