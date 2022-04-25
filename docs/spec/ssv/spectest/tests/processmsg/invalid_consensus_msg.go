package processmsg

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// InvalidConsensusMsg tests an invalid consensus SSVMessage data
func InvalidConsensusMsg() *tests.SpecTest {
	dr := testingutils.BaseRunner()
	startingValue := testingutils.TestAttesterConsensusDataByts
	dr.NewExecutionState()
	if err := dr.StartNewConsensusInstance(startingValue); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.MessageIDForValidatorPKAndRole(testingutils.TestingValidatorPubKey[:], beacon.RoleTypeAttester),
			Data:    []byte{1, 2, 3, 4},
		},
	}

	return &tests.SpecTest{
		Name:                    "ssv msg invalid consensus data",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "c4eb0bb42cc382e468b2362e9d9cc622f388eef6a266901535bb1dfcc51e8868",
		ExpectedError:           "could not get consensus Message from network Message: invalid character '\\x01' looking for beginning of value",
	}
}
