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
// @Summary Create a new job
// @Description Create and queue a new background job. Supports scheduled, recurring, and DAG dependency jobs.
// @Tags jobs
// @Accept json
// @Produce json
// @Param job body job.Job true "Job payload"
// @Success 200 {string} string "Job created: {id}"
// @Failure 400 {string} string "Invalid request"
// @Router /jobs [post]
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
// @Description Retrieve the current list of jobs with their statuses
// @Tags jobs
// @Produce json
// @Success 200 {array} job.Job "List of jobs"
// @Router /jobs/list [get]
func ListJobs(st *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs := st.List()
		json.NewEncoder(w).Encode(jobs)
	}
}
