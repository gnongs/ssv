package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv/duty"
	"github.com/bloxapp/ssv/docs/spec/types"
)

// DutyRunners is a map of duty runners mapped by msg id hex.
type DutyRunners map[beacon.RoleType]*duty.Runner

// DutyRunnerForMsgID returns a Runner from the provided msg ID, or nil if not found
func (ci DutyRunners) DutyRunnerForMsgID(msgID types.MessageID) *duty.Runner {
	role := msgID.GetRoleType()
	return ci[role]
}

type Network interface {
	Broadcast(message types.Encoder) error
}

// Storage is a persistent storage for the SSV
type Storage interface {
}

type BeaconNode interface {
	// GetBeaconNetwork returns the beacon network the node is on
	GetBeaconNetwork() BeaconNetwork
	// GetAttestationData returns attestation data by the given slot and committee index
	GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error)
	// SubmitAttestation submit the attestation to the node
	SubmitAttestation(attestation *phase0.Attestation) error
	// GetBeaconBlock returns beacon block by the given slot and committee index
	GetBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (*phase0.BeaconBlock, error)
}
