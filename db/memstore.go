package db

import (
  "sync"
  "time"

  "github.com/google/uuid"
)

type MemoryStore struct {
  byHash map[string]*ImagePair
  pairs []*ImagePair
  lock sync.RWMutex
}

func NewMemoryStore() Store {
  return &MemoryStore{
    byHash: make(map[string]*ImagePair),
  }
}

func (s *MemoryStore) GetLastImages(n int) []*ImagePair {
  s.lock.RLock()
  defer s.lock.RUnlock()

  var l int
  if len(s.pairs) >= n {
    l = n
  } else {
    l = len(s.pairs)
  }

  return s.pairs[:l]
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

  s.pairs = append([]*ImagePair{pair}, s.pairs...)
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
