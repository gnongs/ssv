package testingutils

import (
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var BaseValidator = func() *ssv.Validator {
	ret := ssv.NewValidator(
		NewTestingNetwork(),
		NewTestingBeaconNode(),
		NewTestingStorage(),
		TestingShare,
		NewTestingKeyManager(),
	)
	ret.DutyRunners[types.BNRoleAttester] = AttesterRunner()
	ret.DutyRunners[types.BNRoleProposer] = ProposerRunner()
	ret.DutyRunners[types.BNRoleAggregator] = AggregatorRunner()
	ret.DutyRunners[types.BNRoleSyncCommittee] = SyncCommitteeRunner()
	return ret
}
