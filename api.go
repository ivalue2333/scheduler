package scheduler

import (
	"context"

	"github.com/ivalue2333/scheduler/dal/rdm"
	"gorm.io/gorm"
)

// 任务状态枚举
const (
	TaskStatusWaiting  = 1 // 等待中
	TaskStatusRunning  = 2 // 执行中
	TaskStatusPending  = 3 // 暂停中
	TaskStatusStop     = 4 // 停止
	TaskStatusFinished = 5 // 完成
	TaskStatusFailed   = 6 // 失败
	TaskStatusTimeout  = 7 // 过期
)

// 任务类型枚举
const (
	// TaskTypeSimple 普通任务，被调度到就运行，且运行一次
	TaskTypeSimple = 1
	// TaskTypeRepeat 重复任务，被调度到就运行，且运行后重新等待调度
	TaskTypeRepeat = 2
	// TaskTypeDelay 延迟任务，延迟一段时间后开始运行，没有到时间不运行，延迟任务有一个预计运行时间，这个预计运行时间尽量保证
	TaskTypeDelay = 3
)

const (
	schedulerKeyPurgeLock = "purge_lock"
	defaultQueueCap       = 1024
)

func MustInit(ctx context.Context, db *gorm.DB, namespace string, category string, opts ...Option) {
	mustInit(ctx, db, namespace, category, opts...)
}

// MustRegister register task name map task func.
func MustRegister(ctx context.Context, taskName string, taskFunc TaskFunc) {
	singleScheduler.mustRegister(ctx, taskName, taskFunc)
}

func SubmitTask(ctx context.Context, task *rdm.SchedulerTask) (int64, error) {
	return singleScheduler.submitTask(ctx, task)
}
