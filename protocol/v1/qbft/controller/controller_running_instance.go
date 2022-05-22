package controller

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/protocol/v1/message"
	"github.com/bloxapp/ssv/protocol/v1/qbft"
	"github.com/bloxapp/ssv/protocol/v1/qbft/instance"
	"github.com/bloxapp/ssv/protocol/v1/sync/changeround"
)

// startInstanceWithOptions will start an iBFT instance with the provided options.
// Does not pre-check instance validity and start validity!
func (c *Controller) startInstanceWithOptions(instanceOpts *instance.Options, value []byte) (*instance.InstanceResult, error) {
	c.currentInstance = instance.NewInstance(instanceOpts)
	c.currentInstance.Init()
	stageChan := c.currentInstance.GetStageChan()

	// reset leader seed for sequence
	if err := c.currentInstance.Start(value); err != nil {
		return nil, errors.WithMessage(err, "could not start iBFT instance")
	}

	metricsCurrentSequence.WithLabelValues(c.Identifier.GetRoleType().String(), hex.EncodeToString(c.Identifier.GetValidatorPK())).Set(float64(c.currentInstance.State().GetHeight()))

	// catch up if we can
	go c.fastChangeRoundCatchup(c.currentInstance)

	// main instance callback loop
	var retRes *instance.InstanceResult
	var err error
instanceLoop:
	for {
		stage := <-stageChan
		if c.currentInstance == nil {
			c.logger.Debug("stage channel was invoked but instance is already empty", zap.Any("stage", stage))
			break instanceLoop
		}
		exit, e := c.instanceStageChange(stage)
		if e != nil {
			err = e
			break instanceLoop
		}
		if exit {
			// exited with no error means instance decided
			// fetch decided msg and return
			retMsg, e := c.ibftStorage.GetDecided(c.Identifier, instanceOpts.Height, instanceOpts.Height)
			if e != nil {
				err = e
				break instanceLoop
			}
			if len(retMsg) == 0 {
				err = errors.New("could not fetch decided msg after instance finished")
				break instanceLoop
			}
			retRes = &instance.InstanceResult{
				Decided: true,
				Msg:     retMsg[0],
			}
			break instanceLoop
		}
	}
	var seq message.Height
	if c.currentInstance != nil {
		// saves seq as instance will be cleared
		seq = c.currentInstance.State().GetHeight()
		// when main instance loop breaks, nil current instance
		c.currentInstance = nil
	}
	c.logger.Debug("iBFT instance result loop stopped")

	c.afterInstance(seq, retRes, err)

	return retRes, err
}

// afterInstance is triggered after the instance was finished
func (c *Controller) afterInstance(height message.Height, res *instance.InstanceResult, err error) {
	// if instance was decided -> wait for late commit messages
	decided := res != nil && res.Decided
	if decided && err == nil {
		if height == message.Height(0) {
			if res.Msg == nil || res.Msg.Message == nil {
				// missing sequence number
				return
			}
			height = res.Msg.Message.Height
		}
		return
	}
	// didn't decided -> purge messages
	//c.q.Purge(msgqueue.DefaultMsgIndex(message.SSVConsensusMsgType, c.Identifier)) // TODO: that's the right indexer? might need be height and all messages
	idn := hex.EncodeToString(c.Identifier)
	c.q.Clean(func(s string) bool {
		if strings.Contains(s, idn) && strings.Contains(s, fmt.Sprintf("%d", height)) {
			return true
		}
		return false
	})
}

// instanceStageChange processes a stage change for the current instance, returns true if requires stopping the instance after stage process.
func (c *Controller) instanceStageChange(stage qbft.RoundState) (bool, error) {
	c.logger.Debug("instance stage has been changed!", zap.String("stage", qbft.RoundState_name[int32(stage)]))
	switch stage {
	case qbft.RoundState_Prepare:
		if err := c.ibftStorage.SaveCurrentInstance(c.GetIdentifier(), c.currentInstance.State()); err != nil {
			return true, errors.Wrap(err, "could not save prepare msg to storage")
		}
	case qbft.RoundState_Decided:
		run := func() error {
			agg, err := c.currentInstance.CommittedAggregatedMsg()
			if err != nil {
				return errors.Wrap(err, "could not get aggregated commit msg and save to storage")
			}
			if err = c.strategy.SaveDecided(agg); err != nil {
				return errors.Wrap(err, "could not save highest decided message to storage")
			}

			ssvMsg, err := c.currentInstance.GetCommittedAggSSVMessage()
			c.logger.Debug("broadcasting decided message", zap.Any("msg", ssvMsg))
			if err != nil {
				return errors.Wrap(err, "could not get SSV message aggregated commit msg")
			}
			if err = c.network.Broadcast(ssvMsg); err != nil {
				return errors.Wrap(err, "could not broadcast decided message")
			}
			c.logger.Info("decided current instance", zap.String("identifier", agg.Message.Identifier.String()),
				zap.Uint64("seqNum", uint64(agg.Message.Height)))
			return nil
		}

		err := run()
		// call stop after decided in order to prevent race condition
		c.currentInstance.Stop()
		if err != nil {
			return true, err
		}
		return false, nil
	case qbft.RoundState_ChangeRound:
		// set time for next round change
		c.currentInstance.ResetRoundTimer()
		// broadcast round change
		if err := c.currentInstance.BroadcastChangeRound(); err != nil {
			c.logger.Error("could not broadcast round change message", zap.Error(err))
		}

	case qbft.RoundState_Stopped:
		c.logger.Info("current iBFT instance stopped, nilling currentInstance", zap.Uint64("seqNum", uint64(c.currentInstance.State().GetHeight())))
		return true, nil
	}
	return false, nil
}

// fastChangeRoundCatchup fetches the latest change round (if one exists) from every peer to try and fast sync forward.
// This is an active msg fetching instead of waiting for an incoming msg to be received which can take a while
func (c *Controller) fastChangeRoundCatchup(instance instance.Instancer) {
	count := 0
	f := changeround.NewLastRoundFetcher(c.logger, c.network)

	handler := func(msg *message.SignedMessage) error {
		if c.currentInstance == nil {
			return errors.New("current instance is nil.")
		}
		err := c.currentInstance.ChangeRoundMsgValidationPipeline().Run(msg)
		if err != nil {
			return errors.Wrap(err, "invalid msg")
		}
		encodedMsg, err := msg.Encode()
		if err != nil {
			return errors.Wrap(err, "could not encode msg")
		}
		c.q.Add(&message.SSVMessage{
			MsgType: message.SSVConsensusMsgType, // should be consensus type as it change round msg
			ID:      c.Identifier,
			Data:    encodedMsg,
		})
		count++
		return nil
	}

	err := f.GetChangeRoundMessages(c.Identifier, instance.State().GetHeight(), handler)

	if err != nil {
		c.logger.Warn("failed fast change round catchup", zap.Error(err))
		return
	}

	c.logger.Info("fast change round catchup finished", zap.Int("count", count))
}
