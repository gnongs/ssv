package ssv

import (
	"time"
)

// Available networks.
const (
	// PraterNetwork represents the Prater test network.
	PraterNetwork BeaconNetwork = "prater"

	// MainNetwork represents the main network.
	MainNetwork BeaconNetwork = "mainnet"

	// NowTestNetwork is a simple test network with genesis time always equal to now, meaning now is slot 0
	NowTestNetwork BeaconNetwork = "now_test_network"
)

// BeaconNetwork represents the network.
type BeaconNetwork string

// NetworkFromString returns network from the given string value
func NetworkFromString(n string) BeaconNetwork {
	switch n {
	case string(PraterNetwork):
		return PraterNetwork
	case string(MainNetwork):
		return MainNetwork
	case string(NowTestNetwork):
		return NowTestNetwork
	default:
		return ""
	}
}

// ForkVersion returns the fork version of the network.
func (n BeaconNetwork) ForkVersion() []byte {
	switch n {
	case PraterNetwork:
		return []byte{0x00, 0x00, 0x10, 0x20}
	case MainNetwork:
		return []byte{0, 0, 0, 0}
	case NowTestNetwork:
		return []byte{0x99, 0x99, 0x99, 0x99}
	default:
		return nil
	}
}

// MinGenesisTime returns min genesis time value
func (n BeaconNetwork) MinGenesisTime() uint64 {
	switch n {
	case PraterNetwork:
		return 1616508000
	case MainNetwork:
		return 1606824023
	case NowTestNetwork:
		return uint64(time.Now().Unix())
	default:
		return 0
	}
}

// SlotDurationSec returns slot duration
func (n BeaconNetwork) SlotDurationSec() time.Duration {
	return 12 * time.Second
}

// SlotsPerEpoch returns number of slots per one epoch
func (n BeaconNetwork) SlotsPerEpoch() uint64 {
	return 32
}

// EstimatedCurrentSlot returns the estimation of the current slot
func (n BeaconNetwork) EstimatedCurrentSlot() uint64 {
	return n.EstimatedSlotAtTime(time.Now().Unix())
}

// EstimatedSlotAtTime estimates slot at the given time
func (n BeaconNetwork) EstimatedSlotAtTime(time int64) uint64 {
	genesis := int64(n.MinGenesisTime())
	if time < genesis {
		return 0
	}
	return uint64(time-genesis) / uint64(n.SlotDurationSec().Seconds())
}

// EstimatedCurrentEpoch estimates the current epoch
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#compute_start_slot_at_epoch
func (n BeaconNetwork) EstimatedCurrentEpoch() uint64 {
	return n.EstimatedEpochAtSlot(n.EstimatedCurrentSlot())
}

// EstimatedEpochAtSlot estimates epoch at the given slot
func (n BeaconNetwork) EstimatedEpochAtSlot(slot uint64) uint64 {
	return slot / n.SlotsPerEpoch()
}
