package job

import (
	"encoding/json"
	"time"
)

type Priority int

const (
	PriorityHigh   Priority = 1
	PriorityMedium Priority = 2
	PriorityLow    Priority = 3
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
)

type Job struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Priority     Priority        `json:"priority"`
	Payload      json.RawMessage `json:"payload"`
	ScheduledAt  time.Time       `json:"scheduled_at"`
	Recurring    string          `json:"recurring,omitempty"` // e.g., "every_1_minute"
	Status       Status          `json:"status"`
	RetryCount   int             `json:"retry_count"`
	MaxRetries   int             `json:"max_retries"`
	Dependencies []string        `json:"dependencies,omitempty"` // DAG
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Error        string          `json:"error,omitempty"`
}

type RecurringInterval string

const (
	Every1Minute  RecurringInterval = "every_1_minute"
	Every5Minutes RecurringInterval = "every_5_minutes"
	Every1Hour    RecurringInterval = "every_1_hour"
)

func (j *Job) IsReady() bool {
	return time.Now().After(j.ScheduledAt) || time.Now().Equal(j.ScheduledAt)
}

func (j *Job) NextScheduled() time.Time {
	if j.Recurring == "" {
		return time.Time{}
	}
	// Parse interval
	switch RecurringInterval(j.Recurring) {
	case Every1Minute:
		return time.Now().Add(1 * time.Minute)
	case Every5Minutes:
		return time.Now().Add(5 * time.Minute)
	case Every1Hour:
		return time.Now().Add(1 * time.Hour)
	}
	return time.Now().Add(1 * time.Hour)
}
