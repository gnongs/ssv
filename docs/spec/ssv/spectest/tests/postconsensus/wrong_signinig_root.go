package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongSigningRoot tests an invalid PostConsensusMessage DutySigningRoot
func WrongSigningRoot() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsgWithWrongRoot(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "invalid PostConsensusMessage DutySigningRoot",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "a5757d77504f4ba7f62430d1b961a140ef15e87a316db1b713b69a453b179841",
		ExpectedError:           "partial sig invalid: SignedPostConsensusMessage invalid: DutySigningRoot invalid",
	}
}
