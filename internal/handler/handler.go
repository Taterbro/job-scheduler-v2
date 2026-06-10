package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/taterbro/job-scheduler-v2/internal/job"
)

type JobHandler struct{}

func NewJobHandler() *JobHandler {
	return &JobHandler{}
}

func (h *JobHandler) Execute(j *job.Job) error {
	log.Printf("Executing job type: %s", j.Type)
	switch j.Type {
	case "send_email":
		return h.sendEmail(j)
	case "log":
		return h.logJob(j)
	default:
		return errors.New("unknown job type")
	}
}

func (h *JobHandler) sendEmail(j *job.Job) error {
	var payload struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
	}
	if err := json.Unmarshal(j.Payload, &payload); err != nil {
		return err
	}
	// Simulate processing
	fmt.Printf("Sending email to %s: %s\n", payload.To, payload.Subject)
	time.Sleep(100 * time.Millisecond) // simulate work
	// Mock failure sometimes
	if time.Now().Unix()%5 == 0 {
		return errors.New("simulated email failure")
	}
	return nil
}

func (h *JobHandler) logJob(j *job.Job) error {
	log.Printf("Logging payload: %s", string(j.Payload))
	return nil
}
