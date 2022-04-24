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
	if err := dr.StartNewInstance(startingValue); err != nil {
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
		PostDutyRunnerStateRoot: "039e927a1858548cc411afe6442ee7222d285661058621bdb1e02405c0f344d4",
		ExpectedError:           "could not get consensus Message from network Message: invalid character '\\x01' looking for beginning of value",
	}
}
