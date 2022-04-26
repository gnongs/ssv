package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processPostConsensusSig(dutyRunner *DutyRunner, signedMsg *SignedPartialSignatureMessage) error {
	quorum, err := dutyRunner.ProcessPostConsensusMessage(signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	switch dutyRunner.BeaconRoleType {
	case beacon.RoleTypeAttester:
		att, err := dutyRunner.State.ReconstructAttestationSig(v.share.ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus sig")
		}
		if err := v.beacon.SubmitAttestation(att); err != nil {
			return errors.Wrap(err, "could not submit to beacon chain reconstructed attestation")
		}
	default:
		return errors.Errorf("unknown duty post consensus sig %s", dutyRunner.BeaconRoleType.String())
	}
	return nil
}

func (v *Validator) signPostConsensusMsg(msg *PartialSignatureMessage) (*SignedPartialSignatureMessage, error) {
	signature, err := v.signer.SignRoot(msg, types.PartialSignatureType, v.share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign PartialSignatureMessage for PostConsensusPartialSig")
	}

	return &SignedPartialSignatureMessage{
		Message:   msg,
		Signature: signature,
		Signers:   []types.OperatorID{v.share.OperatorID},
	}, nil
}
