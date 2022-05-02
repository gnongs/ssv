package spectest

import (
	"github.com/bloxapp/ssv/docs/spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/qbft/spectest/tests/commit"
)

var AllTests = []*tests.SpecTest{
	tests.HappyFlow(),
	tests.SevenOperators(),

	commit.SingleCommit(),
	commit.MultiSignerWithOverlap(),
	commit.MultiSignerNoOverlap(),
	commit.Decided(),
	commit.NoPrevAcceptedProposal(),
	commit.WrongHeight(),
	commit.WrongRound(),
	commit.ImparsableCommitData(),
	commit.WrongCommitData(),
	commit.WrongSignature(),
}
