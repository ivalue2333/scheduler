package scheduler

import (
	"context"
	"math"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/ivalue2333/queues"
	"github.com/ivalue2333/scheduler/dal/query"
	"github.com/ivalue2333/scheduler/dal/rdm"
	"gorm.io/gorm"
)

var (
	singleScheduler *Scheduler
)

func TaskLess(iv, jv *rdm.SchedulerTask) bool {
	return iv.Priority >= jv.Priority
}

type Config struct {
	namespace              string
	category               string
	maxRunningTaskInstance int
	recoverFunc            func(ctx context.Context, err interface{})
}

func newConfig(namespace, category string) *Config {
	config := &Config{
		namespace:              namespace,
		category:               category,
		maxRunningTaskInstance: MaxRunningTaskInstance,
		recoverFunc:            func(ctx context.Context, err interface{}) {},
	}
	return config
}

type Monitor struct {
	RunningTask atomic.Int64
}

type Scheduler struct {
	config      *Config
	monitor     *Monitor
	taskFuncMap map[string]TaskFunc
	tasks       *queues.SyncPriorityQueue[*rdm.SchedulerTask]
	instance    *rdm.SchedulerInstance
	logger      Logger
}

func getSchedulerInstanceID() int64 {
	return rand.Int63n(int64(math.MaxInt32))
}

func mustInit(ctx context.Context, db *gorm.DB, namespace string, category string, opts ...Option) {
	s := &Scheduler{
		config:      newConfig(namespace, category),
		taskFuncMap: make(map[string]TaskFunc),
		tasks:       queues.NewSyncPriorityQueue[*rdm.SchedulerTask](0, TaskLess),
		logger:      newDefaultLogger(),
	}
	query.SetDefault(db)
	s.mustRegisterInstance(ctx)
	for _, opt := range opts {
		if opt != nil {
			err := opt.Apply(s)
			if err != nil {
				panic(err)
			}
		}
	}
	singleScheduler = s
}

func (s *Scheduler) mustRegisterInstance(ctx context.Context) {
	var err error
	for i := 0; i < 10; i++ {
		err = s.tryRegisterInstance(ctx)
		if err == nil {
			break
		}
	}
	if err != nil {
		panic(err)
	}
}

func (s *Scheduler) tryRegisterInstance(ctx context.Context) error {
	instance := newSchedulerInstance()
	err := query.SchedulerInstance.WithContext(ctx).Create(instance)
	if err != nil {
		return err
	}
	s.instance = instance
	return nil
}

func newSchedulerInstance() *rdm.SchedulerInstance {
	return &rdm.SchedulerInstance{
		ID:          getSchedulerInstanceID(),
		HeartBeatAt: time.Now(),
		CreatedAt:   time.Now(),
	}
}

func (s *Scheduler) mustRegister(ctx context.Context, name string, f TaskFunc) {
	s.taskFuncMap[name] = f
}

func (s *Scheduler) submitTask(ctx context.Context, task *rdm.SchedulerTask) (int64, error) {
	err := query.SchedulerTask.WithContext(ctx).Create(task)
	if err != nil {
		return 0, err
	}
	s.tasks.Push(task)
	return task.ID, nil
}

func (s *Scheduler) Schedule(ctx context.Context) {
	for {
		for s.tasks.Len() > 0 {
			task := s.tasks.Pop()
			s.AsyncRunTask(ctx, task)
		}
		time.Sleep(time.Second)
	}
}

func (s *Scheduler) SafeGo(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.config.recoverFunc(ctx, err)
			}
		}()
		f()
	}()
}

func (s *Scheduler) AsyncRunTask(ctx context.Context, task *rdm.SchedulerTask) {
	s.SafeGo(ctx, func() {
		s.RunTask(ctx, task)
	})
}

func (s *Scheduler) RunTask(ctx context.Context, task *rdm.SchedulerTask) {
	s.monitor.RunningTask.Add(1)
	defer func() {
		s.monitor.RunningTask.Add(-1)
	}()
	taskFunc := s.taskFuncMap[task.Name]
	err := taskFunc(ctx, task)
	if err != nil {
		s.logger.Error(ctx, "exec func err:%v, name:%s", err, task.Name)
		_, _ = s.FinishTask(ctx, task.ID, TaskStatusFailed)
		return
	}
	_, err = s.FinishTask(ctx, task.ID, TaskStatusFinished)
	if err != nil {
		return
	}
}

// FinishTask task is done.
func (s *Scheduler) FinishTask(ctx context.Context, id int64, status int32) (int64, error) {
	return s.UpdateStatus(ctx, id, TaskStatusRunning, status)
}

func (s *Scheduler) UpdateStatus(ctx context.Context, id int64, oldStatus, status int32) (int64, error) {
	rowsAffected, err := query.SchedulerTask.WithContext(ctx).UpdateStatus(id, oldStatus, status)
	if err != nil {
		s.logger.Error(ctx, "update status err:%v, id:%s status:%d", err, id, status)
		return 0, err
	}
	return rowsAffected, nil
}
