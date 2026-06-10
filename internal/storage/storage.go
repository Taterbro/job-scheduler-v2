package storage

import (
	"math/rand"
	"sync"
	"time"

	"github.com/taterbro/job-scheduler-v2/internal/job"
)

type Storage struct {
	jobs  map[string]*job.Job
	mu    sync.RWMutex
	dlq   map[string]*job.Job
	dlqMu sync.RWMutex
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewStorage() *Storage {
	return &Storage{
		jobs: make(map[string]*job.Job),
		dlq:  make(map[string]*job.Job),
	}
}

func (s *Storage) Create(j *job.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j.ID == "" {
		j.ID = generateID()
	}
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()
	if j.Status == "" {
		j.Status = job.StatusPending
	}
	s.jobs[j.ID] = j
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(b)
}

func (s *Storage) Get(id string) (*job.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	return j, ok
}

func (s *Storage) Update(j *job.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j.UpdatedAt = time.Now()
	s.jobs[j.ID] = j
}

func (s *Storage) List() []*job.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []*job.Job
	for _, j := range s.jobs {
		list = append(list, j)
	}
	return list
}

func (s *Storage) AddToDLQ(j *job.Job) {
	s.dlqMu.Lock()
	defer s.dlqMu.Unlock()
	s.dlq[j.ID] = j
}

func (s *Storage) GetDLQ() map[string]*job.Job {
	s.dlqMu.RLock()
	defer s.dlqMu.RUnlock()
	copyDLQ := make(map[string]*job.Job, len(s.dlq))
	for k, v := range s.dlq {
		copyDLQ[k] = v
	}
	return copyDLQ
}
