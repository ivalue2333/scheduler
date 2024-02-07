package scheduler

import (
	"context"

	"github.com/ivalue2333/scheduler/dal/rdm"
	"gorm.io/gorm"
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
