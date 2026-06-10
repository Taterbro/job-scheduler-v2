package worker

import (
	"log"
	"math/rand"
	"time"

	"github.com/taterbro/job-scheduler-v2/internal/handler"
	"github.com/taterbro/job-scheduler-v2/internal/job"
	"github.com/taterbro/job-scheduler-v2/internal/scheduler"
	"github.com/taterbro/job-scheduler-v2/internal/storage"
)

type Worker struct {
	id      string
	sched   *scheduler.Scheduler
	storage *storage.Storage
	handler *handler.JobHandler
	stopCh  chan struct{}
}

func NewWorker(id string, s *scheduler.Scheduler, st *storage.Storage, h *handler.JobHandler) *Worker {
	return &Worker{
		id:      id,
		sched:   s,
		storage: st,
		handler: h,
		stopCh:  make(chan struct{}),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case <-w.stopCh:
				return
			default:
				j := w.sched.GetNextJob()
				if j != nil {
					w.processJob(j)
				} else {
					time.Sleep(1 * time.Second) // poll
				}
			}
		}
	}()
}

func (w *Worker) processJob(j *job.Job) {
	// Duplicate protection: assume storage lock or version check, simplistic here
	log.Printf("Worker %s processing job %s", w.id, j.ID)

	j.Status = job.StatusProcessing
	w.storage.Update(j)

	// Execute handler
	err := w.handler.Execute(j)
	if err != nil {
		w.handleFailure(j)
	} else {
		w.handleSuccess(j)
	}
}

func (w *Worker) handleSuccess(j *job.Job) {
	j.Status = job.StatusCompleted
	w.storage.Update(j)
	// Handle recurring
	if j.Recurring != "" {
		next := j.NextScheduled()
		newJ := *j // copy
		newJ.ScheduledAt = next
		newJ.Status = job.StatusPending
		newJ.RetryCount = 0
		w.storage.Create(&newJ)
		w.sched.AddJob(&newJ) // re-enqueue if ready
	}
	// Check DAG dependents
}

func (w *Worker) handleFailure(j *job.Job) {
	j.RetryCount++
	if j.RetryCount <= j.MaxRetries {
		// backoff
		backoff := calculateBackoff(j.RetryCount)
		j.ScheduledAt = time.Now().Add(backoff)
		j.Status = job.StatusPending
		w.storage.Update(j)
		// re-enqueue later
	} else {
		j.Status = job.StatusFailed
		w.storage.AddToDLQ(j)
		// alert if threshold
	}
	w.storage.Update(j)
}

func calculateBackoff(attempt int) time.Duration {
	bases := []time.Duration{1 * time.Second, 5 * time.Second, 25 * time.Second}
	base := bases[attempt-1]
	// jitter ±20%
	jitter := time.Duration(rand.Int63n(int64(base)/5)) - base/10
	return base + jitter
}
