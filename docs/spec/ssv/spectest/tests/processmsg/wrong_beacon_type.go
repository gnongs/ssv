package processmsg

import (
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongBeaconType tests an SSVMessage with the wrong beacon type
func WrongBeaconType() *tests.SpecTest {
	dr := testingutils.BaseRunner()
	startingValue := testingutils.TestAttesterConsensusDataByts
	dr.NewExecutionState()
	if err := dr.StartNewConsensusInstance(startingValue); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		{
			MsgType: 100,
			MsgID:   types.MessageIDForValidatorPKAndRole(testingutils.TestingValidatorPubKey[:], 100),
			Data:    []byte{1, 2, 3, 4},
		},
	}

	return &tests.SpecTest{
		Name:                    "ssv msg wrong beacon type in msg id",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "c4eb0bb42cc382e468b2362e9d9cc622f388eef6a266901535bb1dfcc51e8868",
		ExpectedError:           "Message invalid: could not find duty runner for msg ID",
	}
}
