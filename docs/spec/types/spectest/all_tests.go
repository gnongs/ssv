package spectest

import (
	"github.com/bloxapp/ssv/docs/spec/types/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types/spectest/tests/consensusdata"
)

var AllTests = []*tests.EncodingSpecTest{
	consensusdata.Encoding(),
}
