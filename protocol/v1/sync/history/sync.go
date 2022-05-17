package history

//
//import (
//	"context"
//	"fmt"
//	"github.com/bloxapp/ssv/protocol/v1/sync"
//	"time"
//
//	"github.com/pkg/errors"
//	"go.uber.org/zap"
//
//	"github.com/bloxapp/ssv/protocol/v1/message"
//	p2pprotocol "github.com/bloxapp/ssv/protocol/v1/p2p"
//)
//
//// GetLastDecided reads last decided message from store
//type GetLastDecided func(i message.Identifier) (*message.SignedMessage, error)
//
//
//// History takes care for syncing decided history
//type History interface {
//	// SyncDecided syncs decided message with other peers in the network
//	SyncDecided(ctx context.Context, identifier message.Identifier, getLastDecided GetLastDecided, handler DecidedHandler) (*message.SignedMessage, error)
//	// SyncDecidedRange syncs decided messages for the given identifier and range
//	SyncRange(ctx context.Context, identifier message.Identifier, handler DecidedHandler, from, to message.Height, targetPeers ...string) error
//}
//
//// history implements History
//type history struct {
//	logger   *zap.Logger
//	syncer   p2pprotocol.Syncer
//	fullNode bool
//}
//
//// New creates a new instance of History
//func New(logger *zap.Logger, syncer p2pprotocol.Syncer, fullNode bool) History {
//	return &history{
//		logger:   logger,
//		syncer:   syncer,
//		fullNode: fullNode,
//	}
//}
//
//func (h *history) SyncDecided(ctx context.Context, identifier message.Identifier, getLastDecided GetLastDecided, handler DecidedHandler) (*message.SignedMessage, error) {
//	logger := h.logger.With(zap.String("identifier", fmt.Sprintf("%x", identifier)))
//	var err error
//	var remoteMsgs []p2pprotocol.SyncResult
//	retries := 2
//	for retries > 0 && len(remoteMsgs) == 0 {
//		retries--
//		remoteMsgs, err = h.syncer.LastDecided(identifier)
//		if err != nil {
//			return nil, errors.Wrap(err, "could not fetch local highest instance during sync")
//		}
//		if len(remoteMsgs) == 0 {
//			time.Sleep(250 * time.Millisecond)
//		}
//	}
//	if len(remoteMsgs) == 0 {
//		logger.Info("node is synced: remote highest decided not found (V0), assuming 0")
//		return nil, nil
//	}
//
//	localMsg, err := getLastDecided(identifier)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not fetch local highest instance during sync")
//	}
//	var localHeight message.Height
//	if localMsg != nil {
//		localHeight = localMsg.Message.Height
//	}
//
//	highest, height, sender := sync.GetHighest(h.logger, localMsg, remoteMsgs...)
//	if highest == nil {
//		logger.Info("node is synced: remote highest decided not found (V1), assuming 0")
//		return nil, nil
//	}
//
//	if height <= localHeight {
//		logger.Info("node is synced: local is higher or equal to remote")
//		return nil, nil // no need to save the latest
//	}
//	if !h.fullNode {
//		logger.Info("got highest remote, no need to sync history in a non full-node", zap.Int64("highest", int64(height)))
//		return highest, nil // return the remote heights decided in order to save as latest
//	}
//
//	logger.Debug("syncing decided range...", zap.Int64("local height", int64(localHeight)), zap.Int64("remote height", int64(height)))
//	err = h.SyncRange(ctx, identifier, handler, localHeight, height, sender)
//	if err != nil {
//		// in optimistic approach we ignore failures and updates last decided message
//		h.logger.Debug("could not get decided in range, skipping", zap.Error(err),
//			zap.Int64("from", int64(localHeight)), zap.Int64("to", int64(height)))
//	} else {
//		logger.Debug("node is synced: remote highest found", zap.Int64("height", int64(height)))
//	}
//	return highest, nil
//}
//
//func (h *history) SyncRange(ctx context.Context, identifier message.Identifier, handler DecidedHandler, from, to message.Height, targetPeers ...string) error {
//	visited := make(map[message.Height]bool)
//	msgs, err := h.syncer.GetHistory(identifier, from, to, targetPeers...)
//	if err != nil {
//		return err
//	}
//
//	for _, msg := range msgs {
//		if ctx.Err() != nil {
//			break
//		}
//		sm, err := sync.ExtractSyncMsg(msg.Msg)
//		if err != nil {
//			h.logger.Warn("failed to extract sync msg", zap.Error(err))
//			continue
//		}
//
//	signedMsgLoop:
//		for _, signedMsg := range sm.Data {
//			height := signedMsg.Message.Height
//			if err := handler(signedMsg); err != nil {
//				h.logger.Warn("could not save decided", zap.Error(err), zap.Int64("height", int64(height)))
//				continue
//			}
//			if visited[height] {
//				continue signedMsgLoop
//			}
//			visited[height] = true
//		}
//	}
//	if len(visited) != int(to-from)+1 {
//		h.logger.Warn("not all messages in range", zap.Any("visited", visited), zap.Uint64("to", uint64(to)), zap.Uint64("from", uint64(from)))
//		return errors.Errorf("not all messages in range were saved (%d out of %d)", len(visited), int(to-from))
//	}
//	return nil
//}
