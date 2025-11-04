//go:generate go run gen/main.go
package manager

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
)

const (
	SlowRetryInterval = 200 * time.Millisecond
	FastRetryInterval = 0 * time.Millisecond
	JitterRange       = 100 * time.Millisecond
)

// TODO Manager on init should LoadWorklows and LoadSteps to init the
// CachedWorkflows.

func NewManager[T any](codec transport.Codec[core.Result, core.Cmd[T]],
	outcomeSer OutcomeSerializer,
	maxCmdsCount int,
	ops ...SetOption,
) (manager *Manager[T], err error) {
	o := Options{
		SlowRetryInterval: SlowRetryInterval,
		FastRetryInterval: FastRetryInterval,
		JitterRange:       JitterRange,
	}
	Apply(ops, &o)

	manager = &Manager[T]{
		outcomeSer: outcomeSer,
		queue:      make(chan CmdRecord[T], maxCmdsCount),
		options:    o,
		mu:         sync.Mutex{},
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.partSeq, err = loadFromStorage(manager, codec, outcomeSer)
	return
}

type Manager[T any] struct {
	outcomeSer OutcomeSerializer
	queue      chan CmdRecord[T]
	suspended  bool
	fastRetry  bool
	options    Options
	partSeq    PartSeq

	mu sync.Mutex
}

func (s *Manager[T]) SaveCmd(at time.Time, bs []byte) (id WorkflowID, err error) {
	var (
		envelope = CmdEnvelope{At: at, Bs: bs}
		ebs      = make([]byte, CmdEnvelopeDTS.Size(envelope))
	)
	CmdEnvelopeDTS.Marshal(envelope, ebs)

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.suspended {
		err = ErrServiceSuspended
		return
	}
	partID, partSeq, err := s.options.Storage.Save(ebs)
	if err != nil {
		return
	}
	if s.partSeq != partSeq-1 {
		err = ErrAnotherOrchestratorDetected
		return
	}
	s.partSeq = partSeq
	id = WorkflowID{PartID: partID, PartSeq: partSeq}
	return
}

func (s *Manager[T]) SaveOutcome(serID ServiceID, workID WorkflowID,
	outcome Outcome,
) (err error) {
	bs, err := s.outcomeSer.Marshal(outcome)
	if err != nil {
		return
	}
	envelope := OutcomeEnvelope{serID, workID, bs}
	bs = make([]byte, OutcomeEnvelopeDTS.Size(envelope))
	OutcomeEnvelopeDTS.Marshal(envelope, bs)
	partID, partSeq, err := s.options.Storage.Save(bs)
	if err != nil {
		return
	}
	if s.partSeq != partSeq-1 {
		err = ErrAnotherOrchestratorDetected
		return
	}
	s.partSeq = partSeq
	return s.options.Storage.SetCompleted(partID, partSeq)
}

func (s *Manager[T]) CompleteCmd(workID WorkflowID) (err error) {
	err = s.options.Storage.SetCompleted(workID.PartID, workID.PartSeq)
	if err == nil {
		s.mu.Lock()
		if s.suspended {
			if len(s.queue) > 0 {
				s.fastRetry = true
			} else {
				s.suspended = false
				s.fastRetry = false
			}
		}
		s.mu.Unlock()
	}
	return
}

func (s *Manager[T]) DelayCmd(record CmdRecord[T]) {
	queueCmdRecord(record, s.queue)
}

func (s *Manager[T]) Suspend() {
	s.mu.Lock()
	s.suspended = true
	s.fastRetry = false
	s.mu.Unlock()
}

func (s *Manager[T]) NextDelayedCmd(ctx context.Context) (record CmdRecord[T],
	err error,
) {
	select {
	case <-ctx.Done():
		err = context.Canceled
		return
	case record = <-s.queue:
		select {
		case <-ctx.Done():
			return CmdRecord[T]{}, context.Canceled
		case <-time.After(s.dalayedCmdLag()):
			return
		}
	}
}

func (s *Manager[T]) dalayedCmdLag() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	lag := s.options.SlowRetryInterval
	if s.fastRetry {
		lag = s.options.FastRetryInterval + s.jitter()
	}
	return lag
}

func (s *Manager[T]) jitter() time.Duration {
	if s.options.JitterRange <= 0 {
		return 0
	}
	// rand.Int63n returns 0 <= n < max, so shift to -jitterRange .. +jitterRange
	delta := rand.Int63n(int64(s.options.JitterRange)*2) - int64(s.options.JitterRange)
	return time.Duration(delta)
}

func queueCmdRecord[T any](cmdRecord CmdRecord[T], queue chan<- CmdRecord[T]) {
	select {
	case queue <- cmdRecord:
	default:
		// TODO
		panic("queue is full")
	}
}
