package testingutils

import (
	"encoding/hex"
	altair "github.com/attestantio/go-eth2-client/spec/altair"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/prysmaticlabs/go-bitfield"
)

var TestingAttestationData = &spec.AttestationData{
	Slot:            12,
	Index:           3,
	BeaconBlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Source: &spec.Checkpoint{
		Epoch: 0,
		Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
	Target: &spec.Checkpoint{
		Epoch: 1,
		Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
}
var TestingAttestationRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f") //[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}
var TestingDuty = func(keySet *TestKeySet) *beacon.Duty {
	pk := spec.BLSPubKey{}
	copy(pk[:], keySet.PK.Serialize()[:])

	return &beacon.Duty{
		Type:                    beacon.RoleTypeAttester,
		PubKey:                  pk,
		Slot:                    12,
		ValidatorIndex:          1,
		CommitteeIndex:          22,
		CommitteesAtSlot:        36,
		ValidatorCommitteeIndex: 11,
	}
}

type testingBeaconNode struct {
}

func NewTestingBeaconNode() *testingBeaconNode {
	return &testingBeaconNode{}
}

// GetBeaconNetwork returns the beacon network the node is on
func (bn *testingBeaconNode) GetBeaconNetwork() ssv.BeaconNetwork {
	return ssv.NowTestNetwork
}

// GetAttestationData returns attestation data by the given slot and committee index
func (bn *testingBeaconNode) GetAttestationData(slot spec.Slot, committeeIndex spec.CommitteeIndex) (*spec.AttestationData, error) {
	return TestingAttestationData, nil
}

// SubmitAttestation submit the attestation to the node
func (bn *testingBeaconNode) SubmitAttestation(attestation *spec.Attestation) error {
	return nil
}

// GetBeaconBlock returns beacon block by the given slot and committee index
func (bn *testingBeaconNode) GetBeaconBlock(slot spec.Slot, committeeIndex spec.CommitteeIndex, graffiti, randao []byte) (*altair.BeaconBlock, error) {
	return TestingBeaconBlock, nil
}

// SubmitBeaconBlock submit the block to the node
func (bn *testingBeaconNode) SubmitBeaconBlock(block *altair.SignedBeaconBlock) error {
	return nil
}

// SubmitAggregateSelectionProof returns an AggregateAndProof object
func (bn *testingBeaconNode) SubmitAggregateSelectionProof(slot spec.Slot, committeeIndex spec.CommitteeIndex, slotSig []byte) (*spec.AggregateAndProof, error) {
	return TestingAggregateAndProof, nil
}

// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
func (bn *testingBeaconNode) SubmitSignedAggregateSelectionProof(msg *spec.SignedAggregateAndProof) error {
	return nil
}

// GetSyncMessageBlockRoot returns beacon block root for sync committee
func (bn *testingBeaconNode) GetSyncMessageBlockRoot() (spec.Root, error) {
	return TestingSyncCommitteeBlockRoot, nil
}

// SubmitSyncMessage submits a signed sync committee msg
func (bn *testingBeaconNode) SubmitSyncMessage(msg *altair.SyncCommitteeMessage) error {
	return nil
}
