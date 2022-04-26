package duty_test

import (
	"github.com/bloxapp/ssv/docs/spec/ssv/duty"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDutyExecutionState_Marshaling(t *testing.T) {
	es := &duty.State{
		RunningInstance: testingutils.BaseInstance(),
	}

	byts, err := es.Encode()
	require.NoError(t, err)

	decoded := &duty.State{}
	require.NoError(t, decoded.Decode(byts))
}
