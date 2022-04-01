package postconsensus

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
		PostDutyRunnerStateRoot: "854de580f7c23c94607d671315b57afef2ca8494859ee5f3d4af235ba50c55bb",
		ExpectedError:           "partial sig invalid: SignedPostConsensusMessage invalid: no SignedPostConsensusMessage signers",
	}
}
