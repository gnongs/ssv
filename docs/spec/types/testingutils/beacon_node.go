package testingutils

import (
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
)

var TestingAttestationData = &spec.AttestationData{
	Slot:            1,
	Index:           3,
	BeaconBlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Source: &spec.Checkpoint{
		Epoch: 1,
		Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
	Target: &spec.Checkpoint{
		Epoch: 2,
		Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
}
var TestingAttestationRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f") //[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}
var TestingValidatorPubKey = func() spec.BLSPubKey {
	// sk - 3515c7d08e5affd729e9579f7588d30f2342ee6f6a9334acf006345262162c6f
	pk, _ := hex.DecodeString("8e80066551a81b318258709edaf7dd1f63cd686a0e4db8b29bbb7acfe65608677af5a527d9448ee47835485e02b50bc0")
	blsPK := spec.BLSPubKey{}
	copy(blsPK[:], pk)
	return blsPK
}()
var TestingDuty = &beacon.Duty{
	Type:                    beacon.RoleTypeAttester,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          1,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	ValidatorCommitteeIndex: 11,
}

type testingBeaconNode struct {
}

func NewTestingBeaconNode() *testingBeaconNode {
	return &testingBeaconNode{}
}

// GetAttestationData returns attestation data by the given slot and committee index
func (bn *testingBeaconNode) GetAttestationData(slot spec.Slot, committeeIndex spec.CommitteeIndex) (*spec.AttestationData, error) {
	return TestingAttestationData, nil
}

// SubmitAttestation submit the attestation to the node
func (bn *testingBeaconNode) SubmitAttestation(attestation *spec.Attestation) error {
	return nil
}
