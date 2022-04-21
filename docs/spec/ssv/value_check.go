package ssv

import "github.com/bloxapp/ssv/docs/spec/qbft"

func BeaconAttestationValueCheck(network BeaconNetwork) qbft.ProposedValueCheck {
	// TODO - check for far future singning? https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/validator.md#protection-best-practices
	// TODO - check attestation target and source for weird values (for example target in future epoch etc)
	// TODO - check duty slot equal to attestation slot
	// TODO - check duty committee index to attestation committee index

	return func(data []byte) error {
		return nil
	}
}
