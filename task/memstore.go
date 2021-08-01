package task

import "sync"

type MemoryTaskStore struct {
  tasks map[string]task
  lock sync.RWMutex
}

func NewMemoryTaskStore() TaskStore {
  return &MemoryTaskStore{make(map[string]task), sync.RWMutex{}}
}

func (s *MemoryTaskStore) Save(task *task) error {
  s.lock.Lock()
  defer s.lock.Unlock()

  s.tasks[task.ID] = *task
  return nil
}

func (s *MemoryTaskStore) Get(id string) (*task, bool) {
  s.lock.RLock()
  defer s.lock.RUnlock()

  task, ok := s.tasks[id]
  return &task, ok
}
