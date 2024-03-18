package scheduler

import (
	"context"

	"github.com/ivalue2333/scheduler/dal/rdm"
)

func TaskLess(iv, jv *rdm.SchedulerTask) bool {
	return iv.Priority >= jv.Priority
}

type TaskFunc func(ctx context.Context, task *rdm.SchedulerTask) error

type Option interface {
	Apply(scheduler *Scheduler) error
}
