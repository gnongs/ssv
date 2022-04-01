package spectest

import (
	tests2 "github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests/consensus"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests/messages"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests/postconsensus"
)

var AllTests = []*tests2.SpecTest{
	tests2.HappyFullFlow(),

	postconsensus.ValidMessage(),
	postconsensus.InvaliSignature(),
	postconsensus.WrongSigningRoot(),
	postconsensus.WrongBeaconChainSig(),
	postconsensus.FutureConsensusState(),
	postconsensus.PastConsensusState(),
	postconsensus.MsgAfterReconstruction(),
	postconsensus.DuplicateMsg(),

	messages.NoMessageSigners(),
	messages.MultipleSigners(),
	messages.MultipleMessageSigners(),
	messages.NoSigners(),
	messages.WrongMsgID(),
	messages.UnknownMsgType(),
	messages.NoData(),
	messages.NoDutyRunner(),

	consensus.UnknownDuty(),
	consensus.WrongDutyRole(),
	consensus.WrongDutyPubKey(),
	consensus.DecidedValueSlotMismatch(),
	consensus.AttestationDataNil(),
}
