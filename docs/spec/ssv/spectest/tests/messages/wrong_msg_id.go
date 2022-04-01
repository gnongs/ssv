package messages

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongMsgID tests a SSVMessage ID which doesn't belong to the validator
func WrongMsgID() *tests.SpecTest {
	dr := testingutils.BaseRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgWrongID(nil, testingutils.PostConsensusAttestationMsgWithNoMsgSigners(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "wrong SSVMessage ID",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b",
		ExpectedError:           "Message invalid: msg ID doesn't match validator ID",
	}
}
