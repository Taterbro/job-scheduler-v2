package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/taterbro/job-scheduler-v2/internal/job"
	"github.com/taterbro/job-scheduler-v2/internal/scheduler"
	"github.com/taterbro/job-scheduler-v2/internal/storage"
)

// CreateJob godoc
// @Summary Create a new background job
// @Description Creates a new job that will be queued, scheduled, and processed by workers. Supports immediate, scheduled, recurring, and DAG-dependent jobs.
// @Tags jobs
// @Accept json
// @Produce plain
// @Param job body job.Job true "Job creation payload"
// @Success 200 {string} string "Job created: {job-id}"
// @Failure 400 {string} string "Invalid JSON payload"
// @Router /jobs [post]
// @ExampleValue {"type":"send_email","priority":1,"payload":{"to":"user@example.com","subject":"Welcome to our service","body":"Hello there!"},"scheduled_at":"2026-06-10T10:00:00Z","recurring":"every_5_minutes","max_retries":3,"dependencies":["parent-job-uuid-123"]}
func CreateJob(sched *scheduler.Scheduler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var j job.Job
		if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if j.MaxRetries == 0 {
			j.MaxRetries = 3
		}
		sched.AddJob(&j)
		fmt.Fprintf(w, "Job created: %s", j.ID)
	}
}

// ListJobs godoc
// @Summary List all jobs
// @Description Returns all jobs in the system with their current status, useful for live monitoring (polling).
// @Tags jobs
// @Produce json
// @Success 200 {array} job.Job "List of all jobs"
// @Router /jobs/list [get]
// @ExampleValue [
//
//	{
//	  "id": "job-uuid-12345",
//	  "type": "send_email",
//	  "priority": 1,
//	  "payload": {"to":"user@example.com","subject":"Welcome"},
//	  "scheduled_at": "2026-06-10T08:30:00Z",
//	  "recurring": "every_5_minutes",
//	  "status": "pending",
//	  "retry_count": 0,
//	  "max_retries": 3,
//	  "dependencies": ["job-uuid-999"],
//	  "created_at": "2026-06-10T08:22:00Z",
//	  "updated_at": "2026-06-10T08:22:00Z",
//	  "error": ""
//	},
//	{
//	  "id": "job-uuid-67890",
//	  "type": "send_email",
//	  "priority": 2,
//	  "status": "completed",
//	  "retry_count": 0,
//	  "max_retries": 3,
//	  "created_at": "2026-06-10T08:22:00Z",
//	  "updated_at": "2026-06-10T08:22:00Z",
//	  "error": ""
//	}
//
// ]
func ListJobs(st *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs := st.List()
		json.NewEncoder(w).Encode(jobs)
	}
}
