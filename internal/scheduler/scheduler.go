package scheduler

import (
	"sync"
	"time"

	"github.com/taterbro/job-scheduler-v2/internal/job"
	"github.com/taterbro/job-scheduler-v2/internal/queue"
	"github.com/taterbro/job-scheduler-v2/internal/storage"
)

type Scheduler struct {
	pq      *queue.PriorityQueue
	storage *storage.Storage
	mu      sync.Mutex
}

func NewScheduler(s *storage.Storage) *Scheduler {
	return &Scheduler{
		pq:      queue.NewPriorityQueue(s),
		storage: s,
	}
}

func (s *Scheduler) AddJob(j *job.Job) {
	s.storage.Create(j)
	if j.IsReady() && len(j.Dependencies) == 0 {
		s.pq.Enqueue(j)
	}
	// log
}

func (s *Scheduler) GetNextJob() *job.Job {
	// Poll for ready jobs
	j := s.pq.Dequeue()
	if j != nil {
		j.Status = job.StatusProcessing
		s.storage.Update(j)
		return j
	}
	return nil
}

// For scheduled jobs, a ticker can move them to queue when due
func (s *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			s.checkScheduledJobs()
		}
	}()
}

func (s *Scheduler) checkScheduledJobs() {
	// In full impl, scan pending jobs that are due
	jobs := s.storage.List()
	for _, j := range jobs {
		if j.Status == job.StatusPending && j.IsReady() && len(j.Dependencies) == 0 {
			s.pq.Enqueue(j)
		}
	}
}
