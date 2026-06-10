package dag

import (
	"github.com/taterbro/job-scheduler-v2/internal/job"
	"github.com/taterbro/job-scheduler-v2/internal/storage"
)

type Manager struct {
	storage *storage.Storage
}

func NewDAGManager(s *storage.Storage) *Manager {
	return &Manager{storage: s}
}

func (m *Manager) CanRun(j *job.Job) bool {
	for _, depID := range j.Dependencies {
		dep, ok := m.storage.Get(depID)
		if !ok || dep.Status != job.StatusCompleted {
			return false
		}
	}
	return true
}

// When job completes, notify dependents (in full, use reverse deps)
