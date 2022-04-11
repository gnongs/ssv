package ssv

import "github.com/bloxapp/ssv/docs/spec/qbft"

func (v *Validator) processDecidedMsg(dutyRunner *DutyRunner, msg *qbft.DecidedMessage) error {
	return dutyRunner.QBFTController.ProcessDecidedMsg(msg)
}
