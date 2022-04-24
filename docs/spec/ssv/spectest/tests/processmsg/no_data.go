package processmsg

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// NoData tests a SSVMessage with no data
func NoData() *tests.SpecTest {
	dr := testingutils.BaseRunner()
	startingValue := testingutils.TestAttesterConsensusDataByts
	if err := dr.StartNewInstance(startingValue); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.MessageIDForValidatorPKAndRole(testingutils.TestingValidatorPubKey[:], beacon.RoleTypeAttester),
			Data:    nil,
		},
	}

	return &tests.SpecTest{
		Name:                    "ssv msg no data",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "039e927a1858548cc411afe6442ee7222d285661058621bdb1e02405c0f344d4",
		ExpectedError:           "Message invalid: msg data is invalid",
	}
}
