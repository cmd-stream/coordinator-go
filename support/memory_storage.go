package support

import (
	"sync"

	man "github.com/cmd-stream/coordinator-go/manager"
)

type Checkpoint man.PartSeq

func NewMemoryStorage(batcher *Batcher) *MemoryStorage {
	return &MemoryStorage{
		sl: [][]byte{}, batcher: batcher, mu: &sync.Mutex{},
	}
}

type MemoryStorage struct {
	sl      [][]byte
	chpt    Checkpoint
	batcher *Batcher
	mu      *sync.Mutex
}

func (s *MemoryStorage) Save(bs []byte) (partID man.PartID, partSeq man.PartSeq,
	err error,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sl = append(s.sl, bs)
	return man.PartID(0), man.PartSeq(len(s.sl) - 1), nil
}

func (s *MemoryStorage) SetCompleted(partID man.PartID, partSeq man.PartSeq) (
	err error,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chpt := s.batcher.Add(int64(partSeq))
	if chpt != 0 {
		s.chpt = Checkpoint(chpt)
	}
	return
}

func (s *MemoryStorage) LoadUncompleted(callback man.LoadCallback) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.load(s.sl[s.chpt:], callback)
}

func (s *MemoryStorage) Load(callback man.LoadCallback) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.load(s.sl, callback)
}

func (s *MemoryStorage) Checkpoint() Checkpoint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.chpt
}

func (s *MemoryStorage) load(sl [][]byte, callback man.LoadCallback) (err error) {
	for i, bs := range sl {
		err = callback(man.PartID(0), man.PartSeq(i), bs)
		if err != nil {
			return
		}
	}
	return
}
