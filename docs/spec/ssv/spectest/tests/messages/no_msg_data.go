package messages

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// NoData tests a nil SSVMessage data
func NoData() *tests.SpecTest {
	dr := testingutils.AttesterRunner()

	msg := testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsgWithNoMsgSigners(testingutils.TestingSK1, 1, qbft.FirstHeight))
	msg.Data = []byte{}
	msgs := []*types.SSVMessage{
		msg,
	}

	return &tests.SpecTest{
		Name:                    "no SSVMessage data",
		Runner:                  dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b",
		ExpectedError:           "Message invalid: msg data is invalid",
	}
}
