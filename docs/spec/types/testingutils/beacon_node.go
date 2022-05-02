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
var TestingAttestationRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")

var TestingBeaconBlock = &altair.BeaconBlock{
	Slot:          12,
	ProposerIndex: 10,
	ParentRoot:    spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	StateRoot:     spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Body: &altair.BeaconBlockBody{
		RANDAOReveal: spec.BLSSignature{},
		ETH1Data: &spec.ETH1Data{
			DepositRoot:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			DepositCount: 100,
			BlockHash:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Graffiti:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		ProposerSlashings: []*spec.ProposerSlashing{},
		AttesterSlashings: []*spec.AttesterSlashing{},
		Attestations: []*spec.Attestation{
			{
				AggregationBits: bitfield.NewBitlist(122),
				Data:            TestingAttestationData,
				Signature:       spec.BLSSignature{},
			},
		},
		Deposits:       []*spec.Deposit{},
		VoluntaryExits: []*spec.SignedVoluntaryExit{},
		SyncAggregate: &altair.SyncAggregate{
			SyncCommitteeBits:      bitfield.NewBitvector512(),
			SyncCommitteeSignature: spec.BLSSignature{},
		},
	},
}
var TestingBeaconBlockRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")
var TestingRandaoRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")

var TestingAggregateAndProof = &spec.AggregateAndProof{
	AggregatorIndex: 1,
	SelectionProof:  spec.BLSSignature{},
	Aggregate: &spec.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Signature:       spec.BLSSignature{},
		Data:            TestingAttestationData,
	},
}
var TestingSignedAggregateAndProofRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")
var TestingSelectionProofRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")

var TestingSyncCommitteeBlockRoot = spec.Root{}

var TestingValidatorPubKey = func() spec.BLSPubKey {
	// sk - 5342fd7051ab252e02acc53c765007817b2dc8bab596862e3f8711513b2092b3
	pk, _ := hex.DecodeString("948fb44582ce25336fdb17122eac64fe5a1afc39174ce92d6013becac116766dc5a778c880dd47de7dfff6a0f86ba42c")
	blsPK := spec.BLSPubKey{}
	copy(blsPK[:], pk)
	return blsPK
}()
var TestingWrongValidatorPubKey = func() spec.BLSPubKey {
	pk, _ := hex.DecodeString("948fb44582ce25336fdb17122eac64fe5a1afc39174ce92d6013becac116766dc5a778c880dd47de7dfff6a0f86ba42b")
	blsPK := spec.BLSPubKey{}
	copy(blsPK[:], pk)
	return blsPK
}()

const TestingDutySlot = 12

var TestingAttesterDuty = &types.Duty{
	Type:                    types.BNRoleAttester,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot,
	ValidatorIndex:          1,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingProposerDuty = &types.Duty{
	Type:                    types.BNRoleProposer,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          1,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingAggregatorDuty = &types.Duty{
	Type:                    types.BNRoleAggregator,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          1,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingSyncCommitteeDuty = &types.Duty{
	Type:                    types.BNRoleSyncCommittee,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot,
	ValidatorIndex:          1,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingUnknownDutyType = &types.Duty{
	Type:                    100,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          1,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingWrongDutyPK = &types.Duty{
	Type:                    types.BNRoleAttester,
	PubKey:                  TestingWrongValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          1,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
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
