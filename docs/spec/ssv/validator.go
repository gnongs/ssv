package ssv

import (
	beacon2 "github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv/duty"
	"github.com/bloxapp/ssv/docs/spec/types"
)

// Validator represents an SSV ETH consensus validator share assigned, coordinates duty execution and more.
// Every validator has a validatorID which is validator's public key.
// Each validator has multiple DutyRunners, for each duty type.
type Validator struct {
	DutyRunners DutyRunners
	network     Network
	beacon      BeaconNode
	storage     Storage
	share       *types.Share
	signer      types.KeyManager
}

func NewValidator(
	network Network,
	beacon BeaconNode,
	storage Storage,
	share *types.Share,
	signer types.KeyManager,
) *Validator {
	return &Validator{
		DutyRunners: map[beacon2.RoleType]*duty.Runner{},
		network:     network,
		beacon:      beacon,
		storage:     storage,
		share:       share,
		signer:      signer,
	}
}
