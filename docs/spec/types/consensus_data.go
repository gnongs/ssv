package types

import (
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ConsensusData holds all relevant duty and data Decided on by consensus
type ConsensusData struct {
	Duty                   *Duty
	AttestationData        *phase0.AttestationData
	BlockData              *altair.BeaconBlock
	AggregateAndProof      *phase0.AggregateAndProof
	SyncCommitteeBlockRoot phase0.Root
}

func (cid *ConsensusData) Encode() ([]byte, error) {
	return json.Marshal(cid)
}

func (cid *ConsensusData) Decode(data []byte) error {
	return json.Unmarshal(data, &cid)
}
