package task

import "sync"

type MemoryTaskStore struct {
  tasks map[string]Task
  lock sync.RWMutex
}

func NewMemoryTaskStore() TaskStore {
  return &MemoryTaskStore{make(map[string]Task), sync.RWMutex{}}
}

func (s *MemoryTaskStore) Save(task *Task) error {
  s.lock.Lock()
  defer s.lock.Unlock()

  s.tasks[task.ID] = *task
  return nil
}

func (s *MemoryTaskStore) Get(id string) (*Task, bool) {
  s.lock.RLock()
  defer s.lock.RUnlock()

  task, ok := s.tasks[id]
  return &task, ok
}
