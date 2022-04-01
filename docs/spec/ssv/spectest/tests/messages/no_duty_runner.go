package messages

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// NoDutyRunner tests an SSVMessage ID that doesn't belong to any duty runner
func NoDutyRunner() *tests.SpecTest {
	dr := testingutils.BaseRunner()

	msg := testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsgWithNoMsgSigners(testingutils.TestingSK1, 1, qbft.FirstHeight))
	msg.MsgID = types.MessageIDForValidatorPKAndRole(testingutils.TestingValidatorPubKey[:], beacon.RoleTypeAggregator)
	msgs := []*types.SSVMessage{
		msg,
	}

	return &tests.SpecTest{
		Name:                    "no duty runner found",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b",
		ExpectedError:           "Message invalid: could not find duty runner for msg ID",
	}
}
