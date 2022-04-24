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
	if err := dr.StartNewInstance(startingValue); err != nil {
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
		PostDutyRunnerStateRoot: "039e927a1858548cc411afe6442ee7222d285661058621bdb1e02405c0f344d4",
		ExpectedError:           "Message invalid: msg ID doesn't match validator ID",
	}
}
