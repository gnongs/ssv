package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgContainer_AddIfDoesntExist(t *testing.T) {
	t.Run("same msg and signers", func(t *testing.T) {
		c := &MsgContainer{
			Msgs: map[Round][]*SignedMessage{},
		}

		added, err := c.AddIfDoesntExist(testingSignedMsg)
		require.NoError(t, err)
		require.True(t, added)

		added, err = c.AddIfDoesntExist(testingSignedMsg)
		require.NoError(t, err)
		require.False(t, added)
	})

	t.Run("same msg different signers", func(t *testing.T) {
		c := &MsgContainer{
			Msgs: map[Round][]*SignedMessage{},
		}

		added, err := c.AddIfDoesntExist(testingSignedMsg)
		require.NoError(t, err)
		require.True(t, added)

		added, err = c.AddIfDoesntExist(SignMsg(TestingSK, 2, TestingMessage))
		require.NoError(t, err)
		require.True(t, added)
	})

	t.Run("same msg common signers", func(t *testing.T) {
		c := &MsgContainer{
			Msgs: map[Round][]*SignedMessage{},
		}

		m := testingSignedMsg.DeepCopy()
		m.Signers = []types.OperatorID{1, 2, 3, 4}
		added, err := c.AddIfDoesntExist(m)
		require.NoError(t, err)
		require.True(t, added)

		m = testingSignedMsg.DeepCopy()
		m.Signers = []types.OperatorID{1, 5, 6, 7}
		added, err = c.AddIfDoesntExist(m)
		require.NoError(t, err)
		require.False(t, added)
	})
}

func TestMsgContainer_Marshaling(t *testing.T) {
	c := &MsgContainer{
		Msgs: map[Round][]*SignedMessage{},
	}
	c.Msgs[1] = []*SignedMessage{testingSignedMsg}

	byts, err := c.Encode()
	require.NoError(t, err)

	decoded := &MsgContainer{}
	require.NoError(t, decoded.Decode(byts))

	decodedByts, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, byts, decodedByts)
}
