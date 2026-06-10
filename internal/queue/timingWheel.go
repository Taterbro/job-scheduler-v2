package queue

import (
	"time"

	"github.com/taterbro/job-scheduler-v2/internal/job"
)

// Simple timing wheel for demo
type TimingWheel struct {
	slots    [][]*job.Job
	tick     time.Duration
	current  int
	numSlots int
}

func NewTimingWheel(tick time.Duration, numSlots int) *TimingWheel {
	return &TimingWheel{
		slots:    make([][]*job.Job, numSlots),
		tick:     tick,
		numSlots: numSlots,
	}
}

func (tw *TimingWheel) Add(j *job.Job) {
	// simplistic
	idx := (int(time.Until(j.ScheduledAt)/tw.tick) % tw.numSlots)
	tw.slots[idx] = append(tw.slots[idx], j)
}
