package db

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStore struct {
  byHash map[string]*ImagePair
  pairs []*ImagePair
  lock *sync.RWMutex
}

func NewMemoryStore() Store {
  return &MemoryStore{
    byHash: make(map[string]*ImagePair),
  }
}

func (s *MemoryStore) GetLastImages(n int) []ImagePair {
  result := []ImagePair(nil)
  var start int

  s.lock.RLock()
  defer s.lock.RUnlock()

  if len(s.pairs) < n {
    start = 0
  } else {
    start = len(s.pairs) - n
  }
  end := len(s.pairs)

  for i := start; i < end; i++ {
    result = append(result, *s.pairs[i])
  }

  return result
}

func (s *MemoryStore) AddImage(
  orig, neg, origMime, negMime, hash string,
) *ImagePair {
  pair := &ImagePair{
    ID: uuid.NewString(),
    Original: string(orig),
    Negative: string(neg),
    CreatedAt: time.Now(),
    Hash: hash,
    OrigMime: origMime,
    NegMime: negMime,
  }

  s.lock.Lock()
  defer s.lock.Unlock()

  s.pairs = append(s.pairs, pair)
  s.byHash[pair.Hash] = pair

  return pair
}

func (s *MemoryStore) FindImageByHash(hash string) (*ImagePair, bool) {
  s.lock.RLock()
  defer s.lock.RUnlock()

  if pair, ok := s.byHash[hash]; ok {
    return pair, ok
  }
  return nil, false
}
