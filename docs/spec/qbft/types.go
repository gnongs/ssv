package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"sync"
)

type Round uint64
type Height int64

const (
	NoRound     = 0 // NoRound represents a nil/ zero round
	FirstRound  = 1 // FirstRound value is the first round in any QBFT instance start
	FirstHeight = 0
)

// Network is a collection of funcs for the QBFT Network
type Network interface {
	Broadcast(msg types.Encoder) error
	BroadcastDecided(msg types.Encoder) error
}

type Storage interface {
	// SaveHighestDecided saves (and potentially overrides) the highest Decided for a specific instance
	SaveHighestDecided(signedMsg *SignedMessage) error
}

// ThreadSafeF makes function execution thread safe
type ThreadSafeF struct {
	t sync.Mutex
}

// NewThreadSafeF returns a new instance of NewThreadSafeF
func NewThreadSafeF() *ThreadSafeF {
	return &ThreadSafeF{
		t: sync.Mutex{},
	}
}

// Run runs the provided function
func (safeF *ThreadSafeF) Run(f func() interface{}) interface{} {
	safeF.t.Lock()
	defer safeF.t.Unlock()
	return f()
}
