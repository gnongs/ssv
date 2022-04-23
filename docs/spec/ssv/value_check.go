package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func dutyValueCheck(duty *beacon.Duty, network BeaconNetwork) error {
	if network.EstimatedEpochAtSlot(duty.Slot) > network.EstimatedCurrentEpoch()+1 {
		return errors.New("duty epoch is into far future")
	}
	return nil
}

func BeaconAttestationValueCheck(network BeaconNetwork) qbft.ProposedValueCheck {
	// TODO - check for far future singning? https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/validator.md#protection-best-practices
	// TODO - check attestation target and source for weird values (for example target in future epoch etc)
	// TODO - check duty slot equal to attestation slot
	// TODO - check duty committee index to attestation committee index

	return func(data []byte) error {
		cd := &types.ConsensusData{}
		if err := cd.Decode(data); err != nil {
			return errors.Wrap(err, "failed decoding consensus data")
		}

		if err := dutyValueCheck(cd.Duty, network); err != nil {
			return errors.Wrap(err, "duty invalid")
		}

		if cd.Duty.Type != beacon.RoleTypeAttester {
			return errors.New("duty type != RoleTypeAttester")
		}

		if cd.AttestationData == nil {
			return errors.New("attestation data nil")
		}

		if cd.Duty.Slot != cd.AttestationData.Slot {
			return errors.New("attestation data slot != duty slot")
		}

		if cd.Duty.CommitteeIndex != cd.AttestationData.Index {
			return errors.New("attestation data CommitteeIndex != duty CommitteeIndex")
		}

		// no need to test far future attestation as we check duty slot not far future && duty slot == attestation slot

		if cd.AttestationData.Target.Epoch > network.EstimatedCurrentEpoch()+1 {
			return errors.New("attestation data target epoch is into far future")
		}

		if cd.AttestationData.Source.Epoch >= cd.AttestationData.Target.Epoch {
			return errors.New("attestation data source and target epochs invalid")
		}

		return nil
	}
}
