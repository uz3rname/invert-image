package task

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

type TaskStatus int

const (
  TaskNew TaskStatus = iota
  TaskRunning
  TaskFailed
  TaskDone
  TaskCanceled
)

func (s TaskStatus) String() string {
  switch s {
  case TaskNew:
    return "new"
  case TaskRunning:
    return "running"
  case TaskFailed:
    return "failed"
  case TaskDone:
    return "done"
  case TaskCanceled:
    return "canceled"
  }

  return ""
}

type TaskInput interface{}

type TaskResult interface{}

type TaskHandler interface{
  Run(ctx context.Context, input TaskInput) (TaskResult, error)
}

type Task struct {
  ID string
  Handler TaskHandler
  Status TaskStatus
  Input TaskInput
  Result TaskResult
  Err error
}

func newTask(handler TaskHandler, input TaskInput) *Task {
  return &Task{
    ID: uuid.NewString(),
    Handler: handler,
    Status: TaskNew,
    Input: input,
  }
}

func (t *Task) clearData() {
  empty := &struct{}{}
  t.Input = empty
  t.Result = empty
}

type TaskStore interface {
  Save(task *Task) error
  Get(id string) (*Task, bool)
}

type TaskManager struct {
  handlers map[string]TaskHandler
  taskStore TaskStore
  queue chan *Task
  logger *log.Logger
  cancelFuncs map[string]context.CancelFunc
}

func NewTaskManager(store TaskStore) *TaskManager {
  manager := &TaskManager{
    handlers: make(map[string]TaskHandler),
    taskStore: store,
    queue: make(chan *Task, 1024),
    logger: log.New(os.Stdout, "| ", log.Ltime | log.Lmsgprefix),
    cancelFuncs: make(map[string]context.CancelFunc),
  }
  return manager
}

func (m *TaskManager) log(format string, v ...interface{}) {
  m.logger.Printf(format, v...)
}

func (m *TaskManager) RegisterHandler(
  name string,
  handler TaskHandler,
) error {
  if _, ok := m.handlers[name]; ok {
    return errors.New("Handler '" + name + "' already exists")
  }

  m.handlers[name] = handler
  return nil
}

func (m *TaskManager) PushTask(name string, args TaskInput) (string, error) {
  handler, ok := m.handlers[name]
  if !ok {
    return "", errors.New("Handler '" + name + "' not found")
  }

  task := newTask(handler, args)
  err := m.taskStore.Save(task)
  if err != nil {
    return "", err
  }

  m.queue <-task
  m.log("New task: %s", task.ID)

  return task.ID, nil
}

func (m *TaskManager) RunTask(name string, args TaskInput) (TaskResult, error) {
  handler, ok := m.handlers[name]
  if !ok {
    return "", errors.New("Handler '" + name + "' not found")
  }

  result, err := handler.Run(context.Background(), args)
  return result, err
}

func (m *TaskManager) callCancel(id string) {
  if f, ok := m.cancelFuncs[id]; ok {
    f()
    delete(m.cancelFuncs, id)
  }
}

func (m *TaskManager) CancelTask(id string) error {
  task, ok := m.taskStore.Get(id)
  if !ok {
    return errors.New("Task '" + id + "' not found")
  }
  switch task.Status {
  case TaskNew:
    task.Status = TaskCanceled
    task.clearData()
  case TaskCanceled, TaskFailed, TaskDone:
    return errors.New("Task '" + id + "' is already finished")
  case TaskRunning:
    m.callCancel(task.ID)
    task.Status = TaskCanceled
  }

  return nil
}

func (m *TaskManager) worker(id int) {
  for task := range m.queue {
    if task.Status != TaskNew {
      continue
    }

    m.log("worker %d: starting task '%s'", id, task.ID)

    ctx, cancel := context.WithCancel(context.Background())
    task.Status = TaskRunning
    m.taskStore.Save(task)
    m.cancelFuncs[task.ID] = cancel

    start := time.Now()
    result, err := task.Handler.Run(ctx, task.Input)
    end := time.Now()

    if task.Status != TaskCanceled {
      if err != nil {
        task.Status = TaskFailed
        task.Err = err

        m.log(
          "worker %d: task '%s' failed with error %s", id, task.ID, err.Error(),
        )
      } else {
        task.Status = TaskDone
        task.Result = result
        t := end.Sub(start)

        m.log("worker %d: task '%s' completed within %s", id, task.ID, t.String())
      }
      m.callCancel(task.ID)
    } else {
      m.log("Task '%s' was cancelled", task.ID)
    }
    task.clearData()
    m.taskStore.Save(task)
  }
}

func (m *TaskManager) StartWorkers(n int) {
  for i := 0; i < n; i++ {
    m.log("Starting worker %d of %d", i, n)
    go m.worker(i)
  }
}

func (m *TaskManager) GetStatus(id string) (TaskStatus, error) {
  if task, ok := m.taskStore.Get(id); ok {
    return task.Status, nil
  }
  return 0, errors.New("Task not found")
}

func (m *TaskManager) GetTask(id string) (*Task, bool) {
  task, ok := m.taskStore.Get(id)
  return task, ok
}
