package fullnode

import (
	"context"
	"github.com/bloxapp/ssv/protocol/v1/message"
	p2pprotocol "github.com/bloxapp/ssv/protocol/v1/p2p"
	"github.com/bloxapp/ssv/protocol/v1/qbft/pipelines"
	qbftstorage "github.com/bloxapp/ssv/protocol/v1/qbft/storage"
	"github.com/bloxapp/ssv/protocol/v1/qbft/strategy"
	"github.com/bloxapp/ssv/protocol/v1/sync/history"
	"github.com/bloxapp/ssv/protocol/v1/sync/lastdecided"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type fullNode struct {
	logger         *zap.Logger
	store          qbftstorage.QBFTStore
	decidedFetcher lastdecided.Fetcher
	historySyncer  history.Syncer
}

// NewFullNodeStrategy creates a new instance of fullNode strategy
func NewFullNodeStrategy(logger *zap.Logger, store qbftstorage.QBFTStore, syncer p2pprotocol.Syncer) strategy.Decided {
	return &fullNode{
		logger:         logger.With(zap.String("who", "FullNodeStrategy")),
		store:          store,
		decidedFetcher: lastdecided.NewLastDecidedFetcher(logger.With(zap.String("who", "LastDecidedFetcher")), syncer),
		historySyncer:  history.NewSyncer(logger.With(zap.String("who", "HistorySyncer")), syncer),
	}
}

func (f *fullNode) Sync(ctx context.Context, identifier message.Identifier, pip pipelines.SignedMessagePipeline) (*message.SignedMessage, error) {
	logger := f.logger.With(zap.String("identifier", identifier.String()))
	highest, sender, localHeight, err := f.decidedFetcher.GetLastDecided(ctx, identifier, func(i message.Identifier) (*message.SignedMessage, error) {
		return f.store.GetLastDecided(i)
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not get last decided from peers")
	}
	logger.Debug("highest decided", zap.Int64("local", int64(localHeight)), zap.Any("highest", highest))
	if highest == nil {
		logger.Debug("could not find highest decided from peers")
		return nil, nil
	}
	if highest.Message.Height > localHeight {
		counter := 0
		err := f.historySyncer.SyncRange(ctx, identifier, func(msg *message.SignedMessage) error {
			if err := pip.Run(msg); err != nil {
				return errors.Wrap(err, "invalid msg")
			}
			//f.logger.Debug("saving synced decided", zap.Int64("h", int64(msg.Message.Height)))
			if err := f.store.SaveDecided(msg); err != nil {
				return errors.Wrap(err, "could not save decided msg to storage")
			}
			counter++
			return nil
		}, localHeight, highest.Message.Height, sender)
		if err != nil {
			return nil, errors.Wrap(err, "could not complete sync")
		}
		warnMsg := ""
		if message.Height(counter-1) < highest.Message.Height-localHeight {
			warnMsg = "could not sync all messages in range"
		} else if message.Height(counter-1) > highest.Message.Height-localHeight {
			warnMsg = "got too many messages during synced"
		}
		if len(warnMsg) > 0 {
			logger.Warn(warnMsg,
				zap.Int("actual", counter), zap.Int64("from", int64(localHeight)),
				zap.Int64("to", int64(highest.Message.Height)))
		}
	}
	return highest, nil
}

func (f *fullNode) ValidateHeight(msg *message.SignedMessage) (bool, error) {
	lastDecided, err := f.store.GetLastDecided(msg.Message.Identifier)
	if err != nil {
		return false, errors.Wrap(err, "failed to get last decided")
	}
	if lastDecided == nil {
		return true, nil
	}
	if msg.Message.Height < lastDecided.Message.Height {
		return false, nil
	}
	return true, nil
}

func (f *fullNode) IsMsgKnown(msg *message.SignedMessage) (bool, *message.SignedMessage, error) {
	msgs, err := f.store.GetDecided(msg.Message.Identifier, msg.Message.Height, msg.Message.Height)
	if err == nil && len(msgs) > 0 {
		return true, msgs[0], nil
	}
	return false, nil, err
}

func (f *fullNode) SaveLateCommit(msg *message.SignedMessage) error {
	return f.store.SaveDecided(msg)
}

func (f *fullNode) UpdateDecided(msg *message.SignedMessage) error {
	return f.store.SaveDecided(msg)
}

func (f *fullNode) GetDecided(identifier message.Identifier, heightRange ...message.Height) ([]*message.SignedMessage, error) {
	if len(heightRange) < 2 {
		return nil, errors.New("missing height range")
	}
	return f.store.GetDecided(identifier, heightRange[0], heightRange[1])
}

// SaveDecided in for fullnode saves both decided and last decided
func (f *fullNode) SaveDecided(signedMsgs ...*message.SignedMessage) error {
	if err := f.store.SaveDecided(signedMsgs...); err != nil {
		return errors.Wrap(err, "could not save decided msg to storage")
	}
	// check if we need to save also last decided
	for _, msg := range signedMsgs {
		last, err := f.store.GetLastDecided(msg.Message.Identifier)
		if err != nil {
			return err
		}
		if last != nil && last.Message.Height > msg.Message.Height {
			continue
		}
		if err := f.store.SaveLastDecided(msg); err != nil {
			return err
		}
	}
	return nil
}
