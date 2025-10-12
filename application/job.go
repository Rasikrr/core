package application

import (
	"context"

	"github.com/Rasikrr/core/log"
)

func (a *App) initJobManager(ctx context.Context) error {
	if len(a.jobs) == 0 {
		log.Info(ctx, "no cron jobs to start")
		return nil
	}
	a.jobManager = NewJobManager(a.cronOptions...)
	for _, job := range a.jobs {
		if err := a.jobManager.AddJob(job); err != nil {
			return err
		}
	}
	a.starters.Add(a.jobManager)
	a.closers.Add(a.jobManager)
	log.Info(ctx, "cron jobs initialized", log.Int("job_count", a.jobManager.JobsCount()))
	return nil
}
