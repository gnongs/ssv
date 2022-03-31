package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// MultipleSigners tests >1 SignedPostConsensusMessage Signers
func MultipleSigners() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	noSignerMsg := testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 1, qbft.FirstHeight)
	noSignerMsg.Signers = []types.OperatorID{1, 2}
	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, noSignerMsg),
	}

	return &tests.SpecTest{
		Name:                    ">1 SignedPostConsensusMessage Signers",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "854de580f7c23c94607d671315b57afef2ca8494859ee5f3d4af235ba50c55bb",
		ExpectedError:           "partial sig invalid: SignedPostConsensusMessage allows 1 signer",
	}
}
