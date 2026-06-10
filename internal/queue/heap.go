package queue

import (
	"container/heap"
	"sync"

	"github.com/taterbro/job-scheduler-v2/internal/job"
	"github.com/taterbro/job-scheduler-v2/internal/storage"
)

type JobHeap []*job.Job

func (h JobHeap) Len() int { return len(h) }

func (h JobHeap) Less(i, j int) bool {
	// 1. Priority (lower number higher priority)
	if h[i].Priority != h[j].Priority {
		return h[i].Priority < h[j].Priority
	}
	// 2. Scheduled time
	if !h[i].ScheduledAt.Equal(h[j].ScheduledAt) {
		return h[i].ScheduledAt.Before(h[j].ScheduledAt)
	}
	// 3. Creation time
	return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h JobHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *JobHeap) Push(x interface{}) {
	*h = append(*h, x.(*job.Job))
}

func (h *JobHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type PriorityQueue struct {
	h       JobHeap
	mu      sync.Mutex
	storage *storage.Storage // for marking
}

func NewPriorityQueue(s *storage.Storage) *PriorityQueue {
	return &PriorityQueue{
		h:       make(JobHeap, 0),
		storage: s,
	}
}

func (pq *PriorityQueue) Enqueue(j *job.Job) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(&pq.h, j)
}

func (pq *PriorityQueue) Dequeue() *job.Job {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if len(pq.h) == 0 {
		return nil
	}
	j := heap.Pop(&pq.h).(*job.Job)
	return j
}

func (pq *PriorityQueue) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return len(pq.h)
}
