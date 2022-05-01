package tests

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

type SpecTest struct {
	Name                    string
	Runner                  *ssv.Runner
	Duty                    *beacon.Duty
	Messages                []*types.SSVMessage
	PostDutyRunnerStateRoot string
	ExpectedError           string
}
