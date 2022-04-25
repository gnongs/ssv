package types

import (
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
)

// ConsensusData holds all relevant duty and data Decided on by consensus
type ConsensusData struct {
	Duty            *beacon.Duty
	AttestationData *phase0.AttestationData
	BlockData       *phase0.BeaconBlock
}

func (cid *ConsensusData) Encode() ([]byte, error) {
	return json.Marshal(cid)
}

func (cid *ConsensusData) Decode(data []byte) error {
	return json.Unmarshal(data, &cid)
}
