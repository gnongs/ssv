package ssv_test

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDutyExecutionState_AddPartialSig(t *testing.T) {
	t.Run("add to empty", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{1},
		})

		require.Len(t, s.PostConsensusSignatures, 1)
	})

	t.Run("add multiple", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{1},
		})
		s.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{2},
		})
		s.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{3},
		})

		require.Len(t, s.PostConsensusSignatures, 3)
	})

	t.Run("add duplicate", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{1},
		})
		s.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{1},
		})
		s.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{
			Signers: []types.OperatorID{3},
		})

		require.Len(t, s.PostConsensusSignatures, 2)
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
		dr.DutyExecutionState = &ssv.DutyExecutionState{
			RunningInstance: inst,
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
		dr.DutyExecutionState = &ssv.DutyExecutionState{
			RunningInstance: inst,
			Quorum:          3,
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

	t.Run("Decided, not collected enough sigs but passed PostConsensusSigCollectionSlotTimeout slots", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		inst := testingutils.BaseInstance()
		inst.State.Decided = true
		dr.DutyExecutionState = &ssv.DutyExecutionState{
			RunningInstance: inst,
			Quorum:          3,
			DecidedValue: &types.ConsensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12 + ssv.PostConsensusSigCollectionSlotTimeout,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.EqualError(t, err, "post consensus sig collection is running")
	})

	t.Run("Decided, not collected enough sigs but passed > PostConsensusSigCollectionSlotTimeout slots", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		}
		inst := testingutils.BaseInstance()
		inst.State.Decided = true
		dr.DutyExecutionState = &ssv.DutyExecutionState{
			RunningInstance: inst,
			Quorum:          3,
			DecidedValue: &types.ConsensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12 + ssv.PostConsensusSigCollectionSlotTimeout + 1,
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
		dr.DutyExecutionState = &ssv.DutyExecutionState{
			PostConsensusSignatures: make(map[types.OperatorID][]byte),
			RunningInstance:         inst,
			Quorum:                  3,
			DecidedValue: &types.ConsensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		dr.DutyExecutionState.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{1}, PartialSignature: []byte{1, 2, 3, 4}})
		dr.DutyExecutionState.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{2}, PartialSignature: []byte{1, 2, 3, 4}})
		dr.DutyExecutionState.AddPostConsensusPartialSig(&ssv.PartialSignatureMessage{Signers: []types.OperatorID{3}, PartialSignature: []byte{1, 2, 3, 4}})
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingutils.TestingValidatorPubKey,
		})
		require.NoError(t, err)
	})
}

func TestDutyRunner_StartNewInstance(t *testing.T) {
	t.Run("value nil", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		require.EqualError(t, dr.StartNewConsensusInstance(nil), "new instance value invalid")
	})

	t.Run("valid start", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		dr.ResetExecutionState()
		require.NoError(t, dr.StartNewConsensusInstance(testingutils.TestAttesterConsensusDataByts))
		require.NotNil(t, dr.DutyExecutionState)
		require.NotNil(t, dr.DutyExecutionState.RunningInstance)
		require.EqualValues(t, 3, dr.DutyExecutionState.Quorum)
	})
}

func TestDutyRunner_PostConsensusStateForHeight(t *testing.T) {
	t.Run("no return", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		require.Nil(t, dr.PostConsensusStateForHeight(10))
	})

	t.Run("returns", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		dr.ResetExecutionState()
		require.NoError(t, dr.StartNewConsensusInstance(testingutils.TestAttesterConsensusDataByts))
		require.NotNil(t, dr.PostConsensusStateForHeight(qbft.FirstHeight))
	})
}

func TestDutyRunner_DecideRunningInstance(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dr := testingutils.BaseRunner()
		dr.DutyExecutionState = &ssv.DutyExecutionState{
			PostConsensusSignatures: make(map[types.OperatorID][]byte),
			Quorum:                  3,
		}
		decidedValue := &types.ConsensusData{
			Duty: &beacon.Duty{
				Type:   beacon.RoleTypeAttester,
				Slot:   12,
				PubKey: testingutils.TestingValidatorPubKey,
			},
			AttestationData: nil,
		}

		require.NoError(t, dr.StartNewConsensusInstance(testingutils.TestAttesterConsensusDataByts))
		_, err := dr.DecideRunningInstance(decidedValue, testingutils.NewTestingKeyManager())
		require.NoError(t, err)
		require.NotNil(t, dr.DutyExecutionState.DecidedValue)
		require.NotNil(t, dr.DutyExecutionState.SignedAttestation)
		require.NotNil(t, dr.DutyExecutionState.PostConsensusSigRoot)
		require.NotNil(t, dr.DutyExecutionState.PostConsensusSignatures)
	})
}
