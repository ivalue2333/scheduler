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
	"gorm.io/gorm/clause"
)

var (
	singleScheduler *Scheduler
	instanceDao     = query.SchedulerInstance
	taskDao         = query.SchedulerTask
	lockerDao       = query.SchedulerLocker
)

type Setting struct {
	maxQueueTask        int64         // queue task max number
	maxRunningTask      int64         // running task max number
	maxInstanceNum      int           // max instance number
	heartbeatPeriod     time.Duration // how often do heartbeat
	taskRunPeriod       time.Duration // how often run task
	taskPurgeLockPeriod time.Duration // how often to get purge lock
	taskPurgePeriod     time.Duration // how often purge task
	taskFetchPeriod     time.Duration // how often fetch task
	heartbeatMultiple   int           // heartbeatPeriod mutiple
}

func newDefaultSetting() *Setting {
	return &Setting{
		maxQueueTask:        256,
		maxRunningTask:      10,
		maxInstanceNum:      10 * 1000,
		heartbeatPeriod:     3 * time.Second,
		taskRunPeriod:       time.Millisecond,
		taskPurgeLockPeriod: 5 * time.Second,
		taskPurgePeriod:     5 * time.Second,
		taskFetchPeriod:     5 * time.Second,
		heartbeatMultiple:   5,
	}
}

type Scheduler struct {
	namespace         string
	category          string
	setting           *Setting
	instance          *rdm.SchedulerInstance
	purgeEnable       bool
	taskFuncMap       map[string]TaskFunc
	tasks             *queues.SyncPriorityQueue[*rdm.SchedulerTask]
	logger            Logger
	runningTaskNumber atomic.Int64
	queueTaskNumber   atomic.Int64
	recoverFunc       func(ctx context.Context, err interface{})
	commonCtxFunc     func() context.Context
	runTaskCtxFunc    func(task *rdm.SchedulerTask) context.Context
}

func getRandomInt64() int64 {
	return rand.Int63n(int64(math.MaxInt32))
}

func mustInit(ctx context.Context, db *gorm.DB, namespace string, category string, opts ...Option) {
	s := &Scheduler{
		namespace:      namespace,
		category:       category,
		setting:        newDefaultSetting(),
		taskFuncMap:    make(map[string]TaskFunc),
		tasks:          queues.NewSyncPriorityQueue[*rdm.SchedulerTask](defaultQueueCap, TaskLess),
		logger:         newDefaultLogger(),
		recoverFunc:    func(ctx context.Context, err interface{}) {},
		commonCtxFunc:  func() context.Context { return context.Background() },
		runTaskCtxFunc: func(task *rdm.SchedulerTask) context.Context { return context.Background() },
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
	// start up scheduler
	singleScheduler.StartUp()
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
	err := instanceDao.WithContext(ctx).Create(instance)
	if err != nil {
		return err
	}
	s.instance = instance
	return nil
}

func (s *Scheduler) getInstanceID() int64 {
	return s.instance.ID
}

func (s *Scheduler) getInstanceTaskPurgeAt() time.Time {
	setting := s.setting
	period := time.Duration(setting.heartbeatMultiple) * setting.heartbeatPeriod
	return time.Now().Add(period)
}

// checkInstanceIsAlive instanceID is alive.
func (s *Scheduler) checkInstanceIsAlive(ctx context.Context, instanceID int64) (bool, error) {
	instance, err := instanceDao.WithContext(ctx).GetInstance(instanceID)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		s.logger.Error(ctx, "get instance err:%v", err)
		return false, err
	}
	odlestHeartbeatAt := s.getInstanceTaskPurgeAt()
	if instance.HeartBeatAt.After(odlestHeartbeatAt) {
		return false, nil
	}
	return true, nil
}

func newSchedulerInstance() *rdm.SchedulerInstance {
	return &rdm.SchedulerInstance{
		ID:          getRandomInt64(),
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
	s.addTaskToQueue(task)
	return task.ID, nil
}

// addTaskToQueue if queueTaskNumber greater than maxQueueTask, just return.
func (s *Scheduler) addTaskToQueue(tasks ...*rdm.SchedulerTask) {
	if s.queueTaskNumber.Load() > s.setting.maxQueueTask {
		return
	}
	s.queueTaskNumber.Add(int64(len(tasks)))
	// sometimes queueTaskNumber will greater than maxQueueTask
	s.tasks.Push(tasks...)
}

func (s *Scheduler) popTaskFromQueue() *rdm.SchedulerTask {
	s.queueTaskNumber.Add(-1)
	return s.tasks.Pop()
}

func (s *Scheduler) StartUp() {
	ctx := context.Background()
	s.SafeGo(ctx, func() {
		s.instanceDoHeartbeat()
	})
	s.SafeGo(ctx, func() {
		s.taskRunSchedule()
	})
	s.SafeGo(ctx, func() {
		s.taskAndInstancePurge()
	})
	s.SafeGo(ctx, func() {
		s.instanceGetPurgeLock()
	})
	s.SafeGo(ctx, func() {
		s.taskFetchFromDB()
	})
}

// instanceDoHeartbeat instance do heart beat
func (s *Scheduler) instanceDoHeartbeat() {
	ticker := time.NewTicker(s.setting.heartbeatPeriod)
	for range ticker.C {
		ctx := s.commonCtxFunc()
		_, err := instanceDao.WithContext(ctx).DoHeartBeat(s.getInstanceID(), time.Now())
		if err != nil {
			s.logger.Error(ctx, "do heartbeat err:%v", err)
		}
	}
}

// ScheduleRun run tasks.
func (s *Scheduler) taskRunSchedule() {
	ticker := time.NewTicker(s.setting.taskRunPeriod)
	for range ticker.C {
		if task := s.takeTask(); task != nil {
			s.AsyncRunTask(task)
		}
	}
}

// taskAndInstancePurge change task to waiting for the task`s instance is dead.
// delete dead instance.
func (s *Scheduler) taskAndInstancePurge() {
	ticker := time.NewTicker(s.setting.taskPurgePeriod)
	for range ticker.C {
		if !s.purgeEnable {
			continue
		}
		ctx := s.commonCtxFunc()
		odlestHeartbeatAt := s.getInstanceTaskPurgeAt()
		instances, err := s.ListAllOfflineInstances(ctx, odlestHeartbeatAt)
		if err != nil {
			continue
		}
		err = s.purgeInstanceAndTask(ctx, instances)
		if err != nil {
			s.logger.Error(ctx, "purge instance task err:%v", err)
			continue
		}
	}
}

// purgeInstanceAndTask change task which instance is offline to waiting.
// we just update status serially for better right rate.
func (s *Scheduler) purgeInstanceAndTask(ctx context.Context, instances []*rdm.SchedulerInstance) (err error) {
	for _, instance := range instances {
		// first purge instance task
		_, innerErr := s.PurgeInstanceTaskByInstance(ctx, instance.ID)
		if innerErr != nil {
			err = innerErr
			continue
		}
		// second purge instance
		innerErr = s.PurgeInstance(ctx, instance.ID)
		if innerErr != nil {
			err = innerErr
		}
	}
	return
}

func newSchedulerLocker(instanceID int64) *rdm.SchedulerLocker {
	return &rdm.SchedulerLocker{
		ID:         getRandomInt64(),
		Key:        schedulerKeyPurgeLock,
		InstanceID: instanceID,
	}
}

func (s *Scheduler) instanceGetPurgeLock() {
	ctx := s.commonCtxFunc()
	ticker := time.NewTicker(s.setting.taskPurgeLockPeriod)
	for range ticker.C {
		s.tryGetPurgeLock(ctx)
	}
}

// tryGetPurgeLock use mysql distribute lock to choose one instance to do purge task.
func (s *Scheduler) tryGetPurgeLock(ctx context.Context) (lock bool) {
	locker, err := query.SchedulerLocker.GetLockInstance(schedulerKeyPurgeLock)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Error(ctx, "get lock instance err:%v", err)
	}
	if err == gorm.ErrRecordNotFound {
		locker = nil
	}
	// current no instance
	if locker == nil {
		locker = newSchedulerLocker(s.instance.ID)
		err = query.SchedulerLocker.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(locker)
		if err != nil {
			s.logger.Error(ctx, "try lock err:%v", err)
		}
		// just return wait for next ticker
		return false
	}

	isAlive, err := s.checkInstanceIsAlive(ctx, locker.InstanceID)
	if err != nil {
		return
	}
	if !isAlive {
		_, err = lockerDao.WithContext(ctx).DeleteLockInstance(schedulerKeyPurgeLock)
		return
	}
	// not equal, not get lock
	if locker.InstanceID != s.instance.ID {
		return false
	}
	s.purgeEnable = true
	return true
}

// taskFetch fetch task from db.
func (s *Scheduler) taskFetchFromDB() {
	ticker := time.NewTicker(s.setting.taskFetchPeriod)
	var id int64
	for range ticker.C {
		if s.tasks.Len() < int(s.setting.maxQueueTask) {
			ctx := s.commonCtxFunc()
			tasks, err := s.ListWaittingTasks(ctx, id)
			if err != nil {
				continue
			}
			id = GetMaxIDByTasks(tasks)
			s.addTaskToQueue(tasks...)
		}
	}
}

func (s *Scheduler) SafeGo(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.recoverFunc(ctx, err)
			}
		}()
		f()
	}()
}

func (s *Scheduler) takeTask() *rdm.SchedulerTask {
	if s.tasks.Len() == 0 {
		return nil
	}
	if s.runningTaskNumber.Load() >= s.setting.maxRunningTask {
		return nil
	}
	return s.popTaskFromQueue()
}

func (s *Scheduler) AsyncRunTask(task *rdm.SchedulerTask) {
	ctx := s.runTaskCtxFunc(task)
	s.SafeGo(ctx, func() {
		s.logger.Debug(ctx, "task is running, task:%v", task)
		s.RunTask(ctx, task)
	})
}

func (s *Scheduler) RunTask(ctx context.Context, task *rdm.SchedulerTask) {
	s.runningTaskNumber.Add(1)
	defer func() {
		s.runningTaskNumber.Add(-1)
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

func (s *Scheduler) ListWaittingTasks(ctx context.Context, id int64) ([]*rdm.SchedulerTask, error) {
	tasks, err := query.SchedulerTask.WithContext(ctx).
		ListWaittingTasks(s.category, s.category, s.getInstanceID(), id)
	if err != nil {
		s.logger.Error(ctx, "list waitting task err:%v, instanceID:%d, id:%d", err, s.getInstanceID(), id)
		return nil, err
	}
	return tasks, err
}

func (s *Scheduler) ListInstanceRunningTasks(ctx context.Context, id int64, instanceID int64) ([]*rdm.SchedulerTask, error) {
	tasks, err := query.SchedulerTask.WithContext(ctx).ListInstanceRunningTasks(id, instanceID)
	if err != nil {
		s.logger.Error(ctx, "list instance running task err:%v", err)
		return nil, err
	}
	return tasks, nil
}

func (s *Scheduler) StartTask(ctx context.Context, ids []int64) (int64, error) {
	return s.UpdateStatus(ctx, ids, TaskStatusWaiting, TaskStatusRunning)
}

// FinishTask task is done.
func (s *Scheduler) FinishTask(ctx context.Context, id int64, status int32) (int64, error) {
	return s.UpdateStatus(ctx, []int64{id}, TaskStatusRunning, status)
}

func (s *Scheduler) UpdateStatus(ctx context.Context, ids []int64, oldStatus, status int32) (int64, error) {
	rowsAffected, err := query.SchedulerTask.WithContext(ctx).UpdateStatus(ids, s.getInstanceID(), oldStatus, status)
	if err != nil {
		s.logger.Error(ctx, "update status err:%v, ids:%s status:%d", err, ids, status)
		return 0, err
	}
	return rowsAffected, nil
}

// PurgeInstance delete instanceID.
func (s *Scheduler) PurgeInstance(ctx context.Context, instanceID int64) error {
	_, err := instanceDao.WithContext(ctx).PurgeInstance(instanceID)
	if err != nil {
		s.logger.Error(ctx, "purge instance err:%v, instanceID:%s", err, instanceID)
		return err
	}
	return nil
}

func (s *Scheduler) PurgeInstanceTaskByInstance(ctx context.Context, instanceID int64) (int64, error) {
	var id, count int64
	for {
		tasks, err := s.ListInstanceRunningTasks(ctx, id, instanceID)
		if err != nil {
			return count, err
		}
		if len(tasks) == 0 {
			break
		}
		taskIDs := GetIDByTasks(tasks)
		rowsAffected, err := s.PurgeInstanceTask(ctx, taskIDs, instanceID)
		if err != nil {
			return 0, err
		}
		count += rowsAffected
		time.Sleep(10 * time.Millisecond)
	}
	return count, nil
}

func GetIDByTasks(tasks []*rdm.SchedulerTask) []int64 {
	res := make([]int64, 0, len(tasks))
	for _, task := range tasks {
		res = append(res, task.ID)
	}
	return res
}

func GetMaxIDByTasks(tasks []*rdm.SchedulerTask) (maxID int64) {
	ids := GetIDByTasks(tasks)
	for _, id := range ids {
		if id > maxID {
			maxID = id
		}
	}
	return maxID
}

func (s *Scheduler) PurgeInstanceTask(ctx context.Context, ids []int64, instanceID int64) (int64, error) {
	rowsAffected, err := query.SchedulerTask.WithContext(ctx).PurgeInstanceTask(ids, instanceID)
	if err != nil {
		s.logger.Error(ctx, "purge task err:%v", err)
		return 0, err
	}
	return rowsAffected, nil
}

func (s *Scheduler) ListAllOfflineInstances(ctx context.Context,
	oldestHeartbeatAt time.Time) ([]*rdm.SchedulerInstance, error) {
	var res []*rdm.SchedulerInstance
	var id int64
	for len(res) < s.setting.maxInstanceNum {
		instances, err := s.ListOfflineInstances(ctx, id, oldestHeartbeatAt)
		if err != nil {
			return nil, err
		}
		if len(instances) == 0 {
			break
		}
		// next cursor
		id = instances[len(instances)-1].ID
		// append data
		res = append(res, instances...)
	}
	return res, nil
}

// ListOfflineInstances list offline instance which is judged by heart_beat_at time.
func (s *Scheduler) ListOfflineInstances(ctx context.Context, id int64,
	oldestHeartbeatAt time.Time) ([]*rdm.SchedulerInstance, error) {
	instances, err := instanceDao.ListOfflineInstances(id, oldestHeartbeatAt)
	if err != nil {
		s.logger.Error(ctx, "list offine instance err:%v", err)
		return nil, err
	}
	return instances, nil
}
