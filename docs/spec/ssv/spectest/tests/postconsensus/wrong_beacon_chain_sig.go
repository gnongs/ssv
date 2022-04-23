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
		PostDutyRunnerStateRoot: "33fd61d17dc89513774a7b566e9dddad28ea5703f83efae63aea69e369c1f367",
		ExpectedError:           "partial sig invalid: could not verify beacon partial Signature: could not verify Signature from iBFT member 1",
	}
}
