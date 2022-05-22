package proto

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

var (
	// ErrDuplicateMsgSigner is thrown when trying to aggregate duplicated signers
	ErrDuplicateMsgSigner = errors.New("can't aggregate messages with similar signers")
)

// Compare returns true if both messages are equal.
// DOES NOT compare signatures
func (msg *Message) Compare(other *Message) bool {
	if msg.Type != other.Type ||
		msg.Round != other.Round ||
		!bytes.Equal(msg.Lambda, other.Lambda) ||
		!bytes.Equal(msg.Value, other.Value) {
		return false
	}

	return true
}

// SigningRoot returns a signing root (bytes)
func (msg *Message) SigningRoot() ([]byte, error) {
	// TODO - consider moving to SSZ
	byts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	_, err = hasher.Write(byts)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

// Sign takes a secret key and signs the Message
func (msg *Message) Sign(sk *bls.SecretKey) (*bls.Sign, error) {
	root, err := msg.SigningRoot()
	if err != nil {
		return nil, err
	}
	return sk.SignByte(root), nil
}

// VerifySig returns true if the justification signed msg verifies against the public key, false if otherwise
func (msg *SignedMessage) VerifySig(pk *bls.PublicKey) (bool, error) {
	return msg.VerifyAggregatedSig([]*bls.PublicKey{pk})
}

// VerifyAggregatedSig returns true if the  signed msg verifies against the public keys, false if otherwise
func (msg *SignedMessage) VerifyAggregatedSig(pks []*bls.PublicKey) (bool, error) {
	if msg.Signature == nil || len(msg.Signature) == 0 {
		return false, errors.New("message signature is invalid")
	}

	if len(pks) == 0 {
		return false, errors.New("pks are invalid")
	}

	// signer uniqueness
	err := verifyUniqueSigners(msg.SignerIds)
	if err != nil {
		return false, err
	}

	root, err := msg.Message.SigningRoot()
	if err != nil {
		return false, err
	}

	// aggregate pks
	var aggPK *bls.PublicKey
	for _, pk := range pks {
		if aggPK == nil {
			aggPK = pk
		} else {
			aggPK.Add(pk)
		}
	}

	sig := &bls.Sign{}
	if err := sig.Deserialize(msg.Signature); err != nil {
		return false, err
	}
	return sig.VerifyByte(aggPK, root), nil
}

// SignersIDString returns all KeyManager's Ids as string
func (msg *SignedMessage) SignersIDString() string {
	ret := ""
	for _, i := range msg.SignerIds {
		ret = fmt.Sprintf("%s, %d", ret, i)
	}
	return ret
}

// Aggregate serialize and aggregates signature and signer ID to signed message
func (msg *SignedMessage) Aggregate(other *SignedMessage) error {
	root, err := msg.Message.SigningRoot()
	if err != nil {
		return err
	}
	otherRoot, err := other.Message.SigningRoot()
	if err != nil {
		return err
	}
	if !bytes.Equal(root, otherRoot) {
		return errors.New("can't aggregate different messages")
	}

	// verify not already aggregated
	for _, id := range msg.SignerIds {
		for _, otherID := range other.SignerIds {
			if id == otherID {
				return ErrDuplicateMsgSigner
			}
		}
	}
	// aggregate
	sig := &bls.Sign{}
	if err := sig.Deserialize(msg.Signature); err != nil {
		return err
	}
	otherSig := &bls.Sign{}
	if err := otherSig.Deserialize(other.Signature); err != nil {
		return err
	}
	sig.Add(otherSig)
	msg.Signature = sig.Serialize()
	msg.SignerIds = append(msg.SignerIds, other.SignerIds...)
	return nil
}

// DeepCopy checks marshalling of SignedMessage and returns it
func (msg *SignedMessage) DeepCopy() (*SignedMessage, error) {
	byts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	ret := &SignedMessage{}
	if err := json.Unmarshal(byts, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Encode returns a msg encoded bytes or error
func (msg *SignedMessage) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *SignedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

// VerifySig returns true if the justification signed msg verifies against the public key, false otherwise
func (d *ChangeRoundData) VerifySig(pk bls.PublicKey) (bool, error) {
	err := verifyUniqueSigners(d.SignerIds)
	if err != nil {
		return false, err
	}
	root, err := d.JustificationMsg.SigningRoot()
	if err != nil {
		return false, err
	}

	sig := bls.Sign{}
	if err := sig.Deserialize(d.JustificationSig); err != nil {
		return false, err
	}

	return sig.VerifyByte(&pk, root), nil
}

func verifyUniqueSigners(singerIds []uint64) error {
	unique := map[uint64]bool{}
	for _, signer := range singerIds {
		if _, found := unique[signer]; !found {
			unique[signer] = true
		} else {
			return errors.New("signers are not unique")
		}
	}
	return nil
}

// AggregateMessages will aggregate given msgs or return error
func AggregateMessages(sigs []*SignedMessage) (*SignedMessage, error) {
	var decided *SignedMessage
	var err error
	for _, msg := range sigs {
		if decided == nil {
			decided, err = msg.DeepCopy()
			if err != nil {
				return nil, errors.Wrap(err, "could not copy message")
			}
		} else {
			if err := decided.Aggregate(msg); err != nil {
				return nil, errors.Wrap(err, "could not aggregate message")
			}
		}
	}

	if decided == nil {
		return nil, errors.New("could not aggregate decided messages, no msgs")
	}

	return decided, nil
}
