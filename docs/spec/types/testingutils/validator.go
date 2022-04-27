package testingutils

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv"
)

var BaseValidator = func() *ssv.Validator {
	ret := ssv.NewValidator(
		NewTestingNetwork(),
		NewTestingBeaconNode(),
		NewTestingStorage(),
		TestingShare,
		NewTestingKeyManager(),
	)
	ret.DutyRunners[beacon.RoleTypeAttester] = AttesterRunner()
	ret.DutyRunners[beacon.RoleTypeProposer] = ProposerRunner()
	return ret
}
