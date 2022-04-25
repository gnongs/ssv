package processmsg

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongPubKey tests an SSVMessage ID with the wrong pubkey
func WrongPubKey() *tests.SpecTest {
	dr := testingutils.BaseRunner()
	startingValue := testingutils.TestAttesterConsensusDataByts
	dr.NewExecutionState()
	if err := dr.StartNewConsensusInstance(startingValue); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		{
			MsgType: 100,
			MsgID:   types.MessageIDForValidatorPKAndRole(testingutils.TestingWrongValidatorPubKey[:], beacon.RoleTypeAttester),
			Data:    []byte{1, 2, 3, 4},
		},
	}

	return &tests.SpecTest{
		Name:                    "ssv msg wrong pubkey in msg id",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "c4eb0bb42cc382e468b2362e9d9cc622f388eef6a266901535bb1dfcc51e8868",
		ExpectedError:           "Message invalid: msg ID doesn't match validator ID",
	}
}
