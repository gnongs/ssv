package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// ProcessPostConsensusMessage process post consensus msg, returns true if it has quorum for partial signatures.
// returns true only once (first time quorum achieved)
func (dr *Runner) ProcessPostConsensusMessage(msg *SignedPartialSignatureMessage) (bool, error) {
	if err := dr.canProcessPostConsensusMsg(msg); err != nil {
		return false, errors.Wrap(err, "can't process post consensus message")
	}

	prevQuorum := dr.State.PostConsensusPartialSig.HasQuorum()

	if err := dr.State.PostConsensusPartialSig.AddSignature(msg.Message); err != nil {
		return false, errors.Wrap(err, "could not add partial post consensus signature")
	}

	if prevQuorum {
		return false, nil
	}

	dr.State.SetFinished()
	return dr.State.PostConsensusPartialSig.HasQuorum(), nil
}

// SignDutyPostConsensus sets the Decided duty and partially signs the Decided data, returns a PartialSignatureMessage to be broadcasted or error
func (dr *Runner) SignDutyPostConsensus(decidedValue *types.ConsensusData, signer types.KeyManager) (*PartialSignatureMessage, error) {
	ret := &PartialSignatureMessage{
		Type:    PostConsensusPartialSig,
		Slot:    decidedValue.Duty.Slot,
		Signers: []types.OperatorID{dr.Share.OperatorID},
	}

	switch dr.BeaconRoleType {
	case beacon.RoleTypeAttester:
		signedAttestation, r, err := signer.SignAttestation(decidedValue.AttestationData, decidedValue.Duty, dr.Share.SharePubKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign attestation")
		}

		dr.State.DecidedValue = decidedValue
		dr.State.SignedAttestation = signedAttestation
		dr.State.PostConsensusPartialSig.SigRoot = r

		ret.SigningRoot = dr.State.PostConsensusPartialSig.SigRoot
		ret.PartialSignature = dr.State.SignedAttestation.Signature[:]
		return ret, nil
	case beacon.RoleTypeProposer:
		signedBlock, r, err := signer.SignBeaconBlock(decidedValue.BlockData, decidedValue.Duty, dr.Share.SharePubKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign block")
		}

		dr.State.DecidedValue = decidedValue
		dr.State.SignedProposal = signedBlock
		dr.State.PostConsensusPartialSig.SigRoot = r

		ret.SigningRoot = dr.State.PostConsensusPartialSig.SigRoot
		ret.PartialSignature = dr.State.SignedProposal.Signature[:]
		return ret, nil
	case beacon.RoleTypeAggregator:
		signed, r, err := signer.SignAggregateAndProof(decidedValue.AggregateAndProof, decidedValue.Duty, dr.Share.SharePubKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign aggregate and proof")
		}

		dr.State.DecidedValue = decidedValue
		dr.State.SignedAggregate = signed
		dr.State.PostConsensusPartialSig.SigRoot = r

		ret.SigningRoot = dr.State.PostConsensusPartialSig.SigRoot
		ret.PartialSignature = dr.State.SignedAggregate.Signature[:]
		return ret, nil
	default:
		return nil, errors.Errorf("unknown duty %s", decidedValue.Duty.Type.String())
	}
}

// canProcessPostConsensusMsg returns true if it can process post consensus message, false if not
func (dr *Runner) canProcessPostConsensusMsg(msg *SignedPartialSignatureMessage) error {
	if dr.State.RunningInstance == nil {
		return errors.New("no running instance")
	}

	if decided, _ := dr.State.RunningInstance.IsDecided(); !decided {
		return errors.New("consensus didn't decide")
	}

	if err := dr.validatePartialSigMsg(msg, dr.State.PostConsensusPartialSig, dr.State.DecidedValue.Duty.Slot); err != nil {
		return errors.Wrap(err, "post consensus msg invalid")
	}

	if dr.postConsensusSigTimeout(dr.BeaconNetwork.EstimatedCurrentSlot()) {
		return errors.New("post consensus sig collection timeout")
	}

	return nil
}

func (dr *Runner) verifyBeaconPartialSignature(msg *PartialSignatureMessage) error {
	if len(msg.Signers) != 1 {
		return errors.New("PartialSignatureMessage allows 1 signer")
	}

	signer := msg.Signers[0]
	signature := msg.PartialSignature
	root := msg.SigningRoot

	for _, n := range dr.Share.Committee {
		if n.GetID() == signer {
			pk := &bls.PublicKey{}
			if err := pk.Deserialize(n.GetPublicKey()); err != nil {
				return errors.Wrap(err, "could not deserialized pk")
			}
			sig := &bls.Sign{}
			if err := sig.Deserialize(signature); err != nil {
				return errors.Wrap(err, "could not deserialized Signature")
			}

			// protect nil root
			root = ensureRoot(root)
			// verify
			if !sig.VerifyByte(pk, root) {
				return errors.Errorf("could not verify Signature from iBFT member %d", signer)
			}
			return nil
		}
	}
	return errors.New("beacon partial Signature signer not found")
}

// postConsensusSigTimeout returns true if collecting post consensus sigs timed out
func (dr *Runner) postConsensusSigTimeout(currentSlot spec.Slot) bool {
	return dr.partialSigCollectionTimeout(dr.State.PostConsensusPartialSig, currentSlot)
}

// ensureRoot ensures that SigningRoot will have sufficient allocated memory
// otherwise we get panic from bls:
// github.com/herumi/bls-eth-go-binary/bls.(*Sign).VerifyByte:738
func ensureRoot(root []byte) []byte {
	n := len(root)
	if n == 0 {
		n = 1
	}
	tmp := make([]byte, n)
	copy(tmp[:], root[:])
	return tmp[:]
}
