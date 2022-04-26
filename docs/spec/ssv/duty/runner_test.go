package duty_test

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/ssv/duty"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDutyExecutionState_AddPartialSig(t *testing.T) {
	t.Run("add to empty", func(t *testing.T) {
		s := ssv.NewTestingDutyExecutionState()
		s.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{1},
		})

		require.Len(t, s.PostConsensusPartialSig.Signatures, 1)
	})

	t.Run("add multiple", func(t *testing.T) {
		s := ssv.NewTestingDutyExecutionState()
		s.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{1},
		})
		s.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{2},
		})
		s.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{3},
		})

		require.Len(t, s.PostConsensusPartialSig.Signatures, 3)
	})

	t.Run("add duplicate", func(t *testing.T) {
		s := ssv.NewTestingDutyExecutionState()
		s.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{1},
		})
		s.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{1},
		})
		s.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{3},
		})

		require.Len(t, s.PostConsensusPartialSig.Signatures, 2)
	})
}

func TestDutyRunner_CanStartNewDuty(t *testing.T) {
	t.Run("no prev start", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type: beacon.RoleTypeAttester,
		})
		require.NoError(t, err)
	})

	t.Run("running instance", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		inst := testingutils.BaseInstance()
		inst.State.Decided = false
		dr.State = &duty.State{
			RunningInstance:         inst,
			PostConsensusPartialSig: duty.NewPartialSigContainer(3),
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.EqualError(t, err, "consensus on duty is running")
	})

	t.Run("Decided but still collecting sigs", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		inst := testingutils.BaseInstance()
		inst.State.Decided = true
		dr.State = &duty.RunnerState{
			RunningInstance:         inst,
			PostConsensusPartialSig: duty.NewPartialSigContainer(3),
			DecidedValue: &types.ConsensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.EqualError(t, err, "post consensus sig collection is running")
	})

	t.Run("Decided, not collected enough sigs but passed DutyExecutionSlotTimeout slots", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		inst := testingutils.BaseInstance()
		inst.State.Decided = true
		dr.State = &duty.RunnerState{
			RunningInstance:         inst,
			PostConsensusPartialSig: duty.NewPartialSigContainer(3),
			DecidedValue: &types.ConsensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12 + duty.DutyExecutionSlotTimeout,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.EqualError(t, err, "post consensus sig collection is running")
	})

	t.Run("Decided, not collected enough sigs but passed > DutyExecutionSlotTimeout slots", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		inst := testingutils.BaseInstance()
		inst.State.Decided = true
		dr.State = &duty.RunnerState{
			RunningInstance:         inst,
			PostConsensusPartialSig: duty.NewPartialSigContainer(3),
			DecidedValue: &types.ConsensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12 + duty.DutyExecutionSlotTimeout + 1,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.NoError(t, err)
	})

	t.Run("Decided, collected enough sigs", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		inst := testingutils.BaseInstance()
		inst.State.Decided = true
		dr.State = &duty.RunnerState{
			RunningInstance:         inst,
			PostConsensusPartialSig: duty.NewPartialSigContainer(3),
			DecidedValue: &types.ConsensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		dr.State.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{1}, PartialSignature: []byte{1, 2, 3, 4}})
		dr.State.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{2}, PartialSignature: []byte{1, 2, 3, 4}})
		dr.State.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{3}, PartialSignature: []byte{1, 2, 3, 4}})
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.NoError(t, err)
	})

	t.Run("proposal not collected enough pre consensus sigs", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		dr.BeaconRoleType = beacon.RoleTypeProposer
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeProposer,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		dr.StartNewDuty(duty)
		dr.State = &duty.RunnerState{
			RandaoPartialSig:        duty.NewPartialSigContainer(3),
			PostConsensusPartialSig: duty.NewPartialSigContainer(3),
		}

		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeProposer,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.EqualError(t, err, "randao consensus sig collection is running")
	})

	t.Run("proposal pre consensus sigs timeout", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		dr.BeaconRoleType = beacon.RoleTypeProposer
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeProposer,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		dr.StartNewDuty(duty)
		dr.State = &duty.RunnerState{
			RandaoPartialSig:        duty.NewPartialSigContainer(3),
			PostConsensusPartialSig: duty.NewPartialSigContainer(3),
		}

		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeProposer,
			Slot:   45,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.NoError(t, err)
	})

	t.Run("pre sigs, decided and collected enough post sigs", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeProposer,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		inst := testingutils.BaseInstance()
		inst.State.Decided = true
		dr.State = &duty.RunnerState{
			RunningInstance:         inst,
			RandaoPartialSig:        duty.NewPartialSigContainer(3),
			PostConsensusPartialSig: duty.NewPartialSigContainer(3),
			DecidedValue: &types.ConsensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}

		// pre
		dr.State.RandaoPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{1}, PartialSignature: []byte{1, 2, 3, 4}})
		dr.State.RandaoPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{2}, PartialSignature: []byte{1, 2, 3, 4}})
		dr.State.RandaoPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{3}, PartialSignature: []byte{1, 2, 3, 4}})

		// post
		dr.State.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{1}, PartialSignature: []byte{1, 2, 3, 4}})
		dr.State.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{2}, PartialSignature: []byte{1, 2, 3, 4}})
		dr.State.PostConsensusPartialSig.AddSignature(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{3}, PartialSignature: []byte{1, 2, 3, 4}})
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeProposer,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.NoError(t, err)
	})

}

func TestDutyRunner_StartNewInstance(t *testing.T) {
	t.Run("value nil", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		require.EqualError(t, dr.Decide(nil), "new instance value invalid")
	})

	t.Run("valid start", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		require.NoError(t, dr.Decide(testingutils.TestAttesterConsensusDataByts))
		require.NotNil(t, dr.State)
		require.NotNil(t, dr.State.RunningInstance)
		require.EqualValues(t, 3, dr.State.PostConsensusPartialSig.Quorum)
	})
}

func TestDutyRunner_PostConsensusStateForHeight(t *testing.T) {
	t.Run("no return", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		require.Nil(t, dr.PostConsensusStateForHeight(10))
	})

	t.Run("returns", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		require.NoError(t, dr.Decide(testingutils.TestAttesterConsensusDataByts))
		require.NotNil(t, dr.PostConsensusStateForHeight(qbft.FirstHeight))
	})
}

func TestDutyRunner_DecideRunningInstance(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		dr.State = ssv.NewTestingDutyExecutionState()
		decidedValue := &types.ConsensusData{
			Duty: &beacon.Duty{
				Type:   beacon.RoleTypeAttester,
				Slot:   12,
				PubKey: testingutils.TestingValidatorPubKey,
			},
			AttestationData: nil,
		}

		require.NoError(t, dr.Decide(testingutils.TestAttesterConsensusDataByts))
		_, err := dr.SignDutyPostConsensus(decidedValue, testingutils.NewTestingKeyManager())
		require.NoError(t, err)
		require.NotNil(t, dr.State.DecidedValue)
		require.NotNil(t, dr.State.SignedAttestation)
		require.NotNil(t, dr.State.PostConsensusPartialSig.SigRoot)
		require.NotNil(t, dr.State.PostConsensusPartialSig.Signatures)
	})
}
