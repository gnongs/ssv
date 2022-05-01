package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// StartDuty starts a duty for the validator
func (v *Validator) StartDuty(duty *beacon.Duty) error {
	dutyRunner := v.DutyRunners[duty.Type]
	if dutyRunner == nil {
		return errors.Errorf("duty type %s not supported", duty.Type.String())
	}

	if err := dutyRunner.StartNewDuty(duty); err != nil {
		return errors.Wrap(err, "can't start new duty")
	}

	switch dutyRunner.BeaconRoleType {
	case beacon.RoleTypeAttester:
		return v.executeAttestationDuty(duty, dutyRunner)
	case beacon.RoleTypeProposer:
		return v.executeBlockProposalDuty(duty, dutyRunner)
	default:
		return errors.Errorf("duty type %s unkwon", duty.Type.String())
	}
	return nil
}

// executeBlockProposalDuty steps:
// 1) sign a partial randao sig and wait for 2f+1 partial sigs from peers
// 2) reconstruct randao and send GetBeaconBlock to BN
// 3) start consensus on duty + block data
// 4) Once consensus decides, sign partial block and broadcast
// 5) collect 2f+1 partial sigs, reconstruct and broadcast valid block sig to the BN
func (v *Validator) executeBlockProposalDuty(duty *beacon.Duty, dutyRunner *Runner) error {
	// sign partial randao
	epoch := v.beacon.GetBeaconNetwork().EstimatedEpochAtSlot(duty.Slot)

	msg, err := dutyRunner.SignRandaoPreConsensus(epoch, duty.Slot, v.signer)
	if err != nil {
		return errors.Wrap(err, "could not sign randao for pre-consensus")
	}

	signature, err := v.signer.SignRoot(msg, types.PartialSignatureType, v.share.SharePubKey)
	if err != nil {
		return errors.Wrap(err, "could not sign PartialSignatureMessage for RandaoPartialSig")
	}
	signedPartialMsg := &SignedPartialSignatureMessage{
		Message:   msg,
		Signature: signature,
		Signers:   []types.OperatorID{v.share.OperatorID},
	}

	// broadcast
	data, err := signedPartialMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode post consensus signature msg")
	}
	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(v.share.ValidatorPubKey, dutyRunner.BeaconRoleType),
		Data:    data,
	}
	if err := v.network.Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial randao sig")
	}
	return nil
}

// executeAttestationDuty steps:
// 1) get attestation data from BN
// 2) start consensus on duty + attestation data
// 3) Once consensus decides, sign partial attestation and broadcast
// 4) collect 2f+1 partial sigs, reconstruct and broadcast valid attestation sig to the BN
func (v *Validator) executeAttestationDuty(duty *beacon.Duty, dutyRunner *Runner) error {
	attData, err := v.beacon.GetAttestationData(duty.Slot, duty.CommitteeIndex)
	if err != nil {
		return errors.Wrap(err, "failed to get attestation data")
	}

	input := &types.ConsensusData{
		Duty:            duty,
		AttestationData: attData,
	}

	if err := dutyRunner.Decide(input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}
