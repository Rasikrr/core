package application

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/Rasikrr/core/interfaces"
	"github.com/Rasikrr/core/log"
	"github.com/robfig/cron/v3"
)

type JobManager struct {
	cron     *cron.Cron
	jobs     []interfaces.Job
	entryIDs []cron.EntryID
	started  atomic.Bool
}

func NewJobManager(options ...cron.Option) *JobManager {
	return &JobManager{
		cron: cron.New(options...),
		jobs: make([]interfaces.Job, 0, 5),
	}
}

// AddJob добавляет джобу в менеджер. Должна быть вызвана до Start().
func (jm *JobManager) AddJob(job interfaces.Job) error {
	if job == nil {
		return errors.New("job cannot be nil")
	}

	if jm.started.Load() {
		return errors.New("cannot add job after manager has started")
	}

	jm.jobs = append(jm.jobs, job)
	return nil
}

// Start запускает все зарегистрированные джобы.
func (jm *JobManager) Start(_ context.Context) error {
	if jm.started.Load() {
		return errors.New("job manager already started")
	}

	entryIDs := make([]cron.EntryID, 0, len(jm.jobs))

	for _, job := range jm.jobs {
		entryID, err := jm.cron.AddJob(job.Schedule(), job)
		if err != nil {
			// Откатываем уже добавленные джобы
			for _, id := range entryIDs {
				jm.cron.Remove(id)
			}
			return fmt.Errorf("failed to add job: %w", err)
		}
		entryIDs = append(entryIDs, entryID)
	}

	jm.entryIDs = entryIDs
	jm.cron.Start()
	jm.started.Store(true)

	return nil
}

// Stop останавливает менеджер и ждёт завершения всех запущенных джоб.
func (jm *JobManager) Close(ctx context.Context) error {
	if !jm.started.Load() {
		return nil
	}

	done := make(chan struct{})

	go func() {
		stopCtx := jm.cron.Stop()
		<-stopCtx.Done()
		close(done)
	}()

	select {
	case <-done:
		jm.started.Store(false)
		log.Info(ctx, "all jobs finished")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// IsRunning возвращает статус менеджера.
func (jm *JobManager) IsRunning() bool {
	return jm.started.Load()
}

// Entries возвращает информацию о всех запланированных джобах.
func (jm *JobManager) Entries() []cron.Entry {
	return jm.cron.Entries()
}

// JobsCount возвращает количество зарегистрированных джоб.
func (jm *JobManager) JobsCount() int {
	return len(jm.jobs)
}
