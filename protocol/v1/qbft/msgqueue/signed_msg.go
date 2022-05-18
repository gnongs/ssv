package msgqueue

import (
	"fmt"
	"github.com/bloxapp/ssv/protocol/v1/message"
	"strconv"
	"strings"
)

// SignedMsgCleaner cleans consensus messages from the queue
// it will clean messages of the given identifier and under the given height
func SignedMsgCleaner(mid message.Identifier, h message.Height) Cleaner {
	return func(k string) bool {
		parts := strings.Split(k, "/")
		if len(parts) < 2 {
			return false // unknown
		}
		parts = parts[1:] // remove empty string
		if parts[0] != message.SSVConsensusMsgType.String() {
			return false
		}
		if parts[2] != mid.String() {
			return false
		}
		if getIndexHeight(parts...) > h {
			return false
		}
		// clean
		return true
	}
}

func signedMsgIndexValidator(msg *message.SSVMessage) *message.SignedMessage {
	if msg == nil {
		return nil
	}
	if msg.MsgType != message.SSVConsensusMsgType && msg.MsgType != message.SSVDecidedMsgType {
		return nil
	}
	sm := &message.SignedMessage{}
	if err := sm.Decode(msg.Data); err != nil {
		return nil
	}
	if sm.Message == nil {
		return nil
	}
	return sm
}

// SignedMsgIndexer is the Indexer used for message.SignedMessage
func SignedMsgIndexer() Indexer {
	return func(msg *message.SSVMessage) string {
		if sm := signedMsgIndexValidator(msg); sm == nil {
			return ""
		} else {
			return SignedMsgIndex(msg.MsgType, msg.ID, sm.Message.Height, sm.Message.MsgType)[0]
		}
	}
}

// SignedMsgIndex indexes a message.SignedMessage by identifier, msg type and height
func SignedMsgIndex(msgType message.MsgType, mid message.Identifier, h message.Height, cmt ...message.ConsensusMessageType) []string {
	var res []string
	for _, mt := range cmt {
		res = append(res, fmt.Sprintf("/%s/id/%s/height/%d/qbft_msg_type/%s", msgType.String(), mid.String(), h, mt.String()))
	}
	return res
}

// DecidedMsgIndexer is the Indexer used for decided message.SignedMessage
func DecidedMsgIndexer() Indexer {
	return func(msg *message.SSVMessage) string {
		if msg.MsgType != message.SSVDecidedMsgType {
			return ""
		}
		if sm := signedMsgIndexValidator(msg); sm == nil {
			return ""
		} else {
			return DecidedMsgIndex(msg.ID)
		}
	}
}

// DecidedMsgIndex indexes a decided message.SignedMessage by identifier, msg type
func DecidedMsgIndex(mid message.Identifier) string {
	return fmt.Sprintf("/%s/id/%s/qbft_msg_type/%s", message.SSVDecidedMsgType.String(), mid.String(), message.CommitMsgType.String())
}

func getIndexHeight(idxParts ...string) message.Height {
	hraw := idxParts[2]
	h, err := strconv.Atoi(hraw)
	if err != nil {
		return 0
	}
	return message.Height(h)
}

// getRound returns the round of the message if applicable
func getRound(msg *message.SSVMessage) (message.Round, bool) {
	sm := message.SignedMessage{}
	if err := sm.Decode(msg.Data); err != nil {
		return 0, false
	}
	if sm.Message == nil {
		return 0, false
	}
	return sm.Message.Round, true
}

// getConsensusMsgType returns the message.ConsensusMessageType of the message if applicable
func getConsensusMsgType(msg *message.SSVMessage) (message.ConsensusMessageType, bool) {
	sm := message.SignedMessage{}
	if err := sm.Decode(msg.Data); err != nil {
		return 0, false
	}
	if sm.Message == nil {
		return 0, false
	}
	return sm.Message.MsgType, true
}
