package spectest

import (
	tests2 "github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests/valcheck/attestations"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests/valcheck/beaconblock"
)

var AllTests = []*tests2.SpecTest{
	//postconsensus.ValidMessage(),
	//postconsensus.InvaliSignature(),
	//postconsensus.WrongSigningRoot(),
	//postconsensus.WrongBeaconChainSig(),
	//postconsensus.FutureConsensusState(),
	//postconsensus.PastConsensusState(),
	//postconsensus.MsgAfterReconstruction(),
	//postconsensus.DuplicateMsg(),
	//
	//messages.NoMessageSigners(),
	//messages.MultipleSigners(),
	//messages.MultipleMessageSigners(),
	//messages.NoSigners(),
	//messages.WrongMsgID(),
	//messages.UnknownMsgType(),
	//messages.NoData(),
	//messages.NoDutyRunner(),
	//
	//valcheck.WrongDutyPubKey(),

	attestations.HappyFlow(),
	//attestations.FarFutureDuty(),
	//attestations.DutySlotNotMatchingAttestationSlot(),
	//attestations.DutyCommitteeIndexNotMatchingAttestations(),
	//attestations.FarFutureAttestationTarget(),
	//attestations.AttestationSourceValid(),
	//attestations.DutyTypeWrong(),
	//attestations.AttestationDataNil(),
	//
	//processmsg.NoData(),
	//processmsg.InvalidConsensusMsg(),
	//processmsg.InvalidDecidedMsg(),
	//processmsg.InvalidPostConsensusMsg(),
	//processmsg.UnknownType(),
	//processmsg.WrongPubKey(),
	//processmsg.WrongBeaconType(),

	beaconblock.HappyFlow(),
}
