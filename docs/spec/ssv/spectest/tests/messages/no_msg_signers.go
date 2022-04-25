package messages

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// NoMessageSigners tests an empty PostConsensusMessage Signers
func NoMessageSigners() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsgWithNoMsgSigners(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "NoSigners PostConsensusMessage",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "cbcefe579470d914c3c230bd45cee06e9c5723460044b278a0c629a742551b02",
		ExpectedError:           "partial post consensus sig invalid: SignedPartialSignatureMessage invalid: invalid PartialSignatureMessage signers",
	}
}
