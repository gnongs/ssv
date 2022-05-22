package p2pv1

import (
	"encoding/hex"
	"github.com/bloxapp/ssv/network"
	forksv1 "github.com/bloxapp/ssv/network/forks/v1"
	"github.com/bloxapp/ssv/protocol/v1/message"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// UseMessageRouter registers a message router to handle incoming messages
func (n *p2pNetwork) UseMessageRouter(router network.MessageRouter) {
	n.msgRouter = router
}

// Peers registers a message router to handle incoming messages
func (n *p2pNetwork) Peers(pk message.ValidatorPK) ([]peer.ID, error) {
	all := make([]peer.ID, 0)
	topics := n.fork.ValidatorTopicID(pk)
	for _, topic := range topics {
		peers, err := n.topicsCtrl.Peers(topic)
		if err != nil {
			return nil, err
		}
		all = append(all, peers...)
	}
	return all, nil
}

// Broadcast publishes the message to all peers in subnet
func (n *p2pNetwork) Broadcast(message message.SSVMessage) error {
	if !n.isReady() {
		return ErrNetworkIsNotReady
	}
	raw, err := n.fork.EncodeNetworkMsg(&message)
	if err != nil {
		return errors.Wrap(err, "could not decode message")
	}
	vpk := message.GetIdentifier().GetValidatorPK()
	topics := n.fork.ValidatorTopicID(vpk)

	for _, topic := range topics {
		if topic == forksv1.UnknownSubnet {
			return errors.New("unknown topic")
		}
		if err := n.topicsCtrl.Broadcast(topic, raw, n.cfg.RequestTimeout); err != nil {
			//return errors.Wrap(err, "could not broadcast message")
			return err
		}
	}
	return nil
}

// Subscribe subscribes to validator subnet
func (n *p2pNetwork) Subscribe(pk message.ValidatorPK) error {
	if !n.isReady() {
		return ErrNetworkIsNotReady
	}

	err := n.subscribe(pk)
	if err == nil {
		n.activeValidatorsLock.Lock()
		pkHex := hex.EncodeToString(pk)
		if !n.activeValidators[pkHex] {
			n.activeValidators[pkHex] = true
		}
		n.activeValidatorsLock.Unlock()
	}

	return nil
}

// subscribe subscribes to validator subnet
func (n *p2pNetwork) subscribe(pk message.ValidatorPK) error {
	topics := n.fork.ValidatorTopicID(pk)
	for _, topic := range topics {
		if topic == forksv1.UnknownSubnet {
			return errors.New("unknown topic")
		}
		if err := n.topicsCtrl.Subscribe(topic); err != nil {
			//return errors.Wrap(err, "could not broadcast message")
			return err
		}
	}

	return nil
}

// Unsubscribe unsubscribes from the validator subnet
func (n *p2pNetwork) Unsubscribe(pk message.ValidatorPK) error {
	if !n.isReady() {
		return ErrNetworkIsNotReady
	}
	topics := n.fork.ValidatorTopicID(pk)
	for _, topic := range topics {
		if topic == forksv1.UnknownSubnet {
			return errors.New("unknown topic")
		}
		if err := n.topicsCtrl.Unsubscribe(topic, false); err != nil {
			//return errors.Wrap(err, "could not broadcast message")
			return err
		}
	}
	n.activeValidatorsLock.Lock()
	pkHex := hex.EncodeToString(pk)
	delete(n.activeValidators, pkHex)
	n.activeValidatorsLock.Unlock()
	return nil
}

// handleIncomingMessages reads messages from the given channel and calls the router, note that this function blocks.
func (n *p2pNetwork) handlePubsubMessages(topic string, msg *pubsub.Message) error {
	if n.msgRouter == nil {
		n.logger.Warn("msg router is not configured")
		return nil
	}
	if msg == nil {
		n.logger.Warn("got nil message", zap.String("topic", topic))
		return nil
	}
	ssvMsg, err := n.fork.DecodeNetworkMsg(msg.GetData())
	if err != nil {
		n.logger.Warn("could not decode message", zap.String("topic", topic), zap.Error(err))
		// TODO: handle..
		return nil
	}
	n.msgRouter.Route(*ssvMsg)
	return nil
}
