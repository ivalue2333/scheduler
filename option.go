package scheduler

import (
	"context"

	"github.com/ivalue2333/scheduler/dal/rdm"
)

type Option interface {
	Apply(scheduler *Scheduler) error
}

type ConfigOption struct {
}

func (c *ConfigOption) Apply(s *Scheduler) error {
	return nil
}

type TaskFunc func(ctx context.Context, task *rdm.SchedulerTask) error
