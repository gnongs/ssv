package discovery

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNsToSubnet(t *testing.T) {
	tests := []struct {
		name     string
		ns       string
		expected int64
		isSubnet bool
	}{
		{
			"dummy string",
			"xxx",
			int64(-1),
			false,
		},
		{
			"invalid int",
			"floodsub:bloxstaking.ssv.xxx",
			int64(-1),
			false,
		},
		{
			"invalid",
			"floodsub:ssv.1",
			int64(-1),
			false,
		},
		{
			"valid",
			"floodsub:bloxstaking.ssv.21",
			int64(21),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.isSubnet, isSubnet(test.ns))
			require.Equal(t, test.expected, nsToSubnet(test.ns))
		})
	}
}
