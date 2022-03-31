package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongBeaconChainSig tests an invalid PostConsensusMessage DutySignature
func WrongBeaconChainSig() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsgWithWrongSig(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "Invalid PostConsensusMessage DutySignature",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "854de580f7c23c94607d671315b57afef2ca8494859ee5f3d4af235ba50c55bb",
		ExpectedError:           "partial sig invalid: could not verify beacon partial Signature: could not verify Signature from iBFT member 1",
	}
}
