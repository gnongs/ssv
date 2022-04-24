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

	if err := dutyRunner.CanStartNewDuty(duty); err != nil {
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
func (v *Validator) executeBlockProposalDuty(duty *beacon.Duty, dutyRunner *DutyRunner) error {

}

// executeAttestationDuty steps:
// 1) get attestation data from BN
// 2) start consensus on duty + attestation data
// 3) Once consensus decides, sign partial attestation and broadcast
// 4) collect 2f+1 partial sigs, reconstruct and broadcast valid attestation sig to the BN
func (v *Validator) executeAttestationDuty(duty *beacon.Duty, dutyRunner *DutyRunner) error {
	attData, err := v.beacon.GetAttestationData(duty.Slot, duty.CommitteeIndex)
	if err != nil {
		return errors.Wrap(err, "failed to get attestation data")
	}

	input := &types.ConsensusData{
		Duty:            duty,
		AttestationData: attData,
	}

	byts, err := input.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode input")
	}

	// validate input
	if err := v.valCheck(byts); err != nil {
		return errors.Wrap(err, "StartDuty input data invalid")
	}

	if err := dutyRunner.StartNewConsensusInstance(byts); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}
