package messages

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// UnknownMsgType tests an unknown SSVMessage type
func UnknownMsgType() *tests.SpecTest {
	dr := testingutils.BaseRunner()

	msg := testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsgWithNoMsgSigners(testingutils.TestingSK1, 1, qbft.FirstHeight))
	msg.MsgType = 3
	msgs := []*types.SSVMessage{
		msg,
	}

	return &tests.SpecTest{
		Name:                    "wrong SSVMessage type",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b",
		ExpectedError:           "PartialSignatureMessage Height doesn't match duty runner's Height",
	}
}
