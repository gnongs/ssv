package messages

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// NoSigners tests an empty SignedPostConsensusMessage Signers
func NoSigners() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	noSignerMsg := testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 1, qbft.FirstHeight)
	noSignerMsg.Signers = []types.OperatorID{}
	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, noSignerMsg),
	}

	return &tests.SpecTest{
		Name:                    "NoSigners SignedPostConsensusMessage",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "a5757d77504f4ba7f62430d1b961a140ef15e87a316db1b713b69a453b179841",
		ExpectedError:           "partial sig invalid: SignedPostConsensusMessage invalid: no SignedPostConsensusMessage signers",
	}
}
