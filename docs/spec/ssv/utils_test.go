package ssv_test

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/duty"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

var testConsensusData = &types.ConsensusData{
	Duty:            testingutils.TestingAttesterDuty,
	AttestationData: testingutils.TestingAttestationData,
}
var TestConsensusDataByts, _ = testConsensusData.Encode()

func NewTestingDutyExecutionState() *duty.State {
	return duty.NewDutyExecutionState(3, qbft.FirstHeight)
}
