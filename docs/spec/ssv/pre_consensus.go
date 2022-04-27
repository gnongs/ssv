package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (dr *Runner) SignRandaoPreConsensus(epoch spec.Epoch, signer types.KeyManager) (*PartialSignatureMessage, error) {
	sig, r, err := signer.SignRandaoReveal(epoch, dr.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign partial randao reveal")
	}

	dr.State.RandaoPartialSig.SigRoot = ensureRoot(r)

	// generate partial sig for randao
	msg := &PartialSignatureMessage{
		Type:             RandaoPartialSig,
		PartialSignature: sig,
		SigningRoot:      ensureRoot(r),
		Signers:          []types.OperatorID{dr.Share.OperatorID},
	}

	return msg, nil
}
