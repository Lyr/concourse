package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/metric"
)

//go:generate counterfeiter . BuildScheduler

type BuildScheduler interface {
	Schedule(
		ctx context.Context,
		logger lager.Logger,
		job db.SchedulerJob,
	) (bool, error)
}

type schedulerRunner struct {
	logger     lager.Logger
	jobFactory db.JobFactory
	scheduler  BuildScheduler

	guardJobScheduling chan struct{}
	running            *sync.Map
}

func NewRunner(logger lager.Logger, jobFactory db.JobFactory, scheduler BuildScheduler, maxJobs uint64) Runner {
	newGuardJobScheduling := make(chan struct{}, maxJobs)
	return &schedulerRunner{
		logger:     logger,
		jobFactory: jobFactory,
		scheduler:  scheduler,

		guardJobScheduling: newGuardJobScheduling,
		running:            &sync.Map{},
	}
}

func (s *schedulerRunner) Run(ctx context.Context) error {
	sLog := s.logger.Session("run")

	sLog.Debug("start")
	defer sLog.Debug("done")

	jobs, err := s.jobFactory.JobsToSchedule()
	if err != nil {
		return fmt.Errorf("find jobs to schedule: %w", err)
	}

	for _, j := range jobs {
		if _, exists := s.running.LoadOrStore(j.ID(), true); exists {
			// already scheduling this job
			continue
		}

		s.guardJobScheduling <- struct{}{}

		jLog := sLog.Session("job", lager.Data{"job": j.Name()})

		go func(job db.SchedulerJob) {
			defer func() {
				<-s.guardJobScheduling
				s.running.Delete(job.ID())
			}()

			schedulingLock, acquired, err := job.AcquireSchedulingLock(sLog)
			if err != nil {
				jLog.Error("failed-to-acquire-lock", err)
				return
			}

			if !acquired {
				return
			}

			defer schedulingLock.Release()

			err = s.scheduleJob(ctx, sLog, job)
			if err != nil {
				jLog.Error("failed-to-schedule-job", err)
			}
		}(j)
	}

	return nil
}

func (s *schedulerRunner) scheduleJob(ctx context.Context, logger lager.Logger, job db.SchedulerJob) error {
	metric.JobsScheduling.Inc()
	defer metric.JobsScheduling.Dec()
	defer metric.JobsScheduled.Inc()

	logger = logger.Session("schedule-job", lager.Data{"job": job.Name()})

	logger.Debug("schedule")

	// Grabs out the requested time that triggered off the job schedule in
	// order to set the last scheduled to the exact time of this triggering
	// request
	requestedTime := job.ScheduleRequestedTime()

	found, err := job.Reload()
	if err != nil {
		return fmt.Errorf("reload job: %w", err)
	}

	if !found {
		logger.Debug("could-not-find-job-to-reload")
		return nil
	}

	jStart := time.Now()

	needsRetry, err := s.scheduler.Schedule(
		ctx,
		logger,
		job,
	)
	if err != nil {
		return fmt.Errorf("schedule job: %w", err)
	}

	if !needsRetry {
		err = job.UpdateLastScheduled(requestedTime)
		if err != nil {
			logger.Error("failed-to-update-last-scheduled", err, lager.Data{"job": job.Name()})
			return fmt.Errorf("update last scheduled: %w", err)
		}
	}

	metric.SchedulingJobDuration{
		PipelineName: job.PipelineName(),
		JobName:      job.Name(),
		JobID:        job.ID(),
		Duration:     time.Since(jStart),
	}.Emit(logger)

	return nil
}
