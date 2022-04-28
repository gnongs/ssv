package types

import (
	"encoding/hex"
	"github.com/bloxapp/ssv/beacon"
	"github.com/stretchr/testify/require"
	"testing"
)

var testingPubKey = make([]byte, 48)

func TestMessageIDForValidatorPKAndRole(t *testing.T) {
	require.EqualValues(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, NewMsgID(testingPubKey, beacon.RoleTypeAttester, 100))
}

func TestMessageID_GetRoleType(t *testing.T) {
	t.Run("attester", func(t *testing.T) {
		msgID := NewMsgID(testingPubKey, beacon.RoleTypeAttester, 100)
		require.EqualValues(t, beacon.RoleTypeAttester, msgID.GetRoleType())
	})

	t.Run("proposer", func(t *testing.T) {
		msgID := NewMsgID(testingPubKey, beacon.RoleTypeProposer, 100)
		require.EqualValues(t, beacon.RoleTypeProposer, msgID.GetRoleType())
	})

	t.Run("long pk", func(t *testing.T) {
		msgID := NewMsgID(testingPubKey, beacon.RoleTypeProposer, 100)
		require.EqualValues(t, beacon.RoleTypeProposer, msgID.GetRoleType())
	})
}

func TestMessageID_GetSlot(t *testing.T) {
	t.Run("100", func(t *testing.T) {
		msgID := NewMsgID(testingPubKey, beacon.RoleTypeAttester, 100)
		require.EqualValues(t, 100, msgID.GetSlot())
	})
	t.Run("1000", func(t *testing.T) {
		msgID := NewMsgID(testingPubKey, beacon.RoleTypeAttester, 1000)
		require.EqualValues(t, 1000, msgID.GetSlot())
	})
	t.Run("100000000000", func(t *testing.T) {
		msgID := NewMsgID(testingPubKey, beacon.RoleTypeAttester, 100000000000)
		require.EqualValues(t, 100000000000, msgID.GetSlot())
	})
}

func TestShare_Marshaling(t *testing.T) {
	expected, _ := hex.DecodeString("7b2264617461223a223232343135313439343434323431336433643232222c226964223a223031303230333034222c2274797065223a223330227d")

	t.Run("encode", func(t *testing.T) {
		msg := &SSVMessage{
			MsgID:   MessageID{1, 2, 3, 4},
			MsgType: SSVConsensusMsgType,
			Data:    []byte{1, 2, 3, 4},
		}

		byts, err := msg.Encode()
		require.NoError(t, err)
		require.EqualValues(t, expected, byts)
	})

	t.Run("decode", func(t *testing.T) {
		msg := &SSVMessage{}
		require.NoError(t, msg.Decode(expected))
		require.EqualValues(t, MessageID{1, 2, 3, 4}, msg.MsgID)
		require.EqualValues(t, SSVConsensusMsgType, msg.MsgType)
		require.EqualValues(t, []byte{1, 2, 3, 4}, msg.Data)
	})
}
