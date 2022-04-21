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
		PostDutyRunnerStateRoot: "a5757d77504f4ba7f62430d1b961a140ef15e87a316db1b713b69a453b179841",
		ExpectedError:           "partial sig invalid: SignedPostConsensusMessage invalid: invalid PostConsensusMessage signers",
	}
}
