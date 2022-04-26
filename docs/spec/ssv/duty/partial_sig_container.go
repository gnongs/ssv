package duty

import (
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

type PartialSigContainer struct {
	Signatures map[types.OperatorID][]byte
	SigRoot    []byte
	// Quorum is the number of min signatures needed for quorum
	Quorum uint64
}

func NewPartialSigContainer(quorum uint64) *PartialSigContainer {
	return &PartialSigContainer{
		Quorum:     quorum,
		Signatures: make(map[types.OperatorID][]byte),
	}
}

func (ps *PartialSigContainer) AddSignature(sigMsg *ssv.PartialSignatureMessage) error {
	if len(sigMsg.Signers) != 1 {
		return errors.New("PartialSignatureMessage has != 1 Signers")
	}

	if ps.Signatures[sigMsg.Signers[0]] == nil {
		ps.Signatures[sigMsg.Signers[0]] = sigMsg.PartialSignature
	}
	return nil
}

func (ps *PartialSigContainer) ReconstructSignature(validatorPubKey []byte) ([]byte, error) {
	// Reconstruct signatures
	signature, err := types.ReconstructSignatures(ps.Signatures)
	if err != nil {
		return nil, errors.Wrap(err, "failed to reconstruct signatures")
	}
	if err := types.VerifyReconstructedSignature(signature, validatorPubKey, ps.SigRoot); err != nil {
		return nil, errors.Wrap(err, "failed to verify reconstruct signature")
	}
	return signature.Serialize(), nil
}

func (ps *PartialSigContainer) HasQuorum() bool {
	return uint64(len(ps.Signatures)) >= ps.Quorum
}
