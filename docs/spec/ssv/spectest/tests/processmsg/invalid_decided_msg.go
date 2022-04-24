package processmsg

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// InvalidDecidedMsg tests an invalid decided SSVMessage data
func InvalidDecidedMsg() *tests.SpecTest {
	dr := testingutils.BaseRunner()
	startingValue := testingutils.TestAttesterConsensusDataByts
	if err := dr.StartNewInstance(startingValue); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVDecidedMsgType,
			MsgID:   types.MessageIDForValidatorPKAndRole(testingutils.TestingValidatorPubKey[:], beacon.RoleTypeAttester),
			Data:    []byte{1, 2, 3, 4},
		},
	}

	return &tests.SpecTest{
		Name:                    "ssv msg invalid decided data",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "039e927a1858548cc411afe6442ee7222d285661058621bdb1e02405c0f344d4",
		ExpectedError:           "could not get decided Message from network Message: invalid character '\\x01' looking for beginning of value",
	}
}
