package qbft_test

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUponCommit(t *testing.T) {
	t.Run("single commit", func(t *testing.T) {
		i := testingutils.BaseInstance()
		i.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
		})
		msg := testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})
		_, _, _, err := qbft.UponCommit(i.State, testingutils.TestingConfig, msg, i.State.CommitContainer)
		require.NoError(t, err)
	})

	t.Run("single commits to quorum", func(t *testing.T) {
		i := testingutils.BaseInstance()
		i.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
		})

		msg := testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})
		decided, _, _, err := qbft.UponCommit(i.State, testingutils.TestingConfig, msg, i.State.CommitContainer)
		require.NoError(t, err)
		require.False(t, decided)

		msg = testingutils.SignQBFTMsg(testingutils.TestingSK2, types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})
		decided, _, _, err = qbft.UponCommit(i.State, testingutils.TestingConfig, msg, i.State.CommitContainer)
		require.NoError(t, err)
		require.False(t, decided)

		msg = testingutils.SignQBFTMsg(testingutils.TestingSK3, types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})
		decided, _, _, err = qbft.UponCommit(i.State, testingutils.TestingConfig, msg, i.State.CommitContainer)
		require.NoError(t, err)
		require.True(t, decided)
	})

	t.Run("multi signer to quorum", func(t *testing.T) {
		i := testingutils.BaseInstance()
		i.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
		})

		msg := testingutils.MultiSignQBFTMsg([]*bls.SecretKey{testingutils.TestingSK1, testingutils.TestingSK2}, []types.OperatorID{1, 2}, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})
		decided, _, _, err := qbft.UponCommit(i.State, testingutils.TestingConfig, msg, i.State.CommitContainer)
		require.NoError(t, err)
		require.False(t, decided)

		msg = testingutils.SignQBFTMsg(testingutils.TestingSK3, types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})
		decided, _, _, err = qbft.UponCommit(i.State, testingutils.TestingConfig, msg, i.State.CommitContainer)
		require.NoError(t, err)
		require.True(t, decided)
	})

	t.Run("multi signer to quorum", func(t *testing.T) {
		i := testingutils.BaseInstance()
		i.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
		})

		msg := testingutils.MultiSignQBFTMsg([]*bls.SecretKey{testingutils.TestingSK1, testingutils.TestingSK2}, []types.OperatorID{1, 2}, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})
		decided, _, _, err := qbft.UponCommit(i.State, testingutils.TestingConfig, msg, i.State.CommitContainer)
		require.NoError(t, err)
		require.False(t, decided)

		msg = testingutils.MultiSignQBFTMsg([]*bls.SecretKey{testingutils.TestingSK3, testingutils.TestingSK4}, []types.OperatorID{3, 4}, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		})
		decided, _, _, err = qbft.UponCommit(i.State, testingutils.TestingConfig, msg, i.State.CommitContainer)
		require.NoError(t, err)
		require.True(t, decided)
	})

}
