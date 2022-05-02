package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/bloxapp/ssv/beacon"
)

// ValidatorPK is an eth2 validator public key
type ValidatorPK []byte

const (
	pubKeySize       = 48
	pubKeyStartPos   = 0
	roleTypeSize     = 4
	roleTypeStartPos = pubKeyStartPos + pubKeySize
)

type Validate interface {
	// Validate returns error if msg validation doesn't pass.
	// Msg validation checks the msg, it's variables for validity.
	Validate() error
}

// MessageIDBelongs returns true if message ID belongs to validator
func (vid ValidatorPK) MessageIDBelongs(msgID MessageID) bool {
	toMatch := msgID.GetPubKey()
	return bytes.Equal(vid, toMatch)
}

// MessageID is used to identify and route messages to the right validator and Runner
type MessageID []byte

func (msg MessageID) GetPubKey() []byte {
	return msg[pubKeyStartPos : pubKeyStartPos+pubKeySize]
}

func (msg MessageID) GetRoleType() beacon.RoleType {
	roleByts := msg[roleTypeStartPos : roleTypeStartPos+roleTypeSize]
	return beacon.RoleType(binary.LittleEndian.Uint32(roleByts))
}

func NewMsgID(pk []byte, role beacon.RoleType) MessageID {
	roleByts := make([]byte, 4)
	binary.LittleEndian.PutUint32(roleByts, uint32(role))
	return append(pk, roleByts...)
}

func (msgID MessageID) String() string {
	return hex.EncodeToString(msgID)
}

type MsgType uint64

const (
	// SSVConsensusMsgType are all QBFT consensus related messages
	SSVConsensusMsgType MsgType = iota
	// SSVDecidedMsgType are all QBFT decided messages
	SSVDecidedMsgType
	// SSVSyncMsgType are all QBFT sync messages
	SSVSyncMsgType
	// SSVPartialSignatureMsgType are all partial signatures msgs over beacon chain specific signatures
	SSVPartialSignatureMsgType
)

type Root interface {
	// GetRoot returns the root used for signing and verification
	GetRoot() ([]byte, error)
}

type MessageSignature interface {
	Root
	GetSignature() Signature
	GetSigners() []OperatorID
	// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
	MatchedSigners(ids []OperatorID) bool
	// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
	Aggregate(signedMsg MessageSignature) error
}
