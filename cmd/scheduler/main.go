package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/taterbro/job-scheduler-v2/internal/handler"
	"github.com/taterbro/job-scheduler-v2/internal/job"
	"github.com/taterbro/job-scheduler-v2/internal/scheduler"
	"github.com/taterbro/job-scheduler-v2/internal/storage"
	"github.com/taterbro/job-scheduler-v2/internal/worker"
)

func main() {
	st := storage.NewStorage()
	sched := scheduler.NewScheduler(st)
	h := handler.NewJobHandler()

	// Start scheduler
	sched.Start()

	// Start workers
	w1 := worker.NewWorker("worker-1", sched, st, h)
	w1.Start()
	// w2 := worker.NewWorker("worker-2", sched, st, h)
	// w2.Start()

	// API for creating jobs
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
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
	})

	// Simple UI or list
	http.HandleFunc("/jobs/list", func(w http.ResponseWriter, r *http.Request) {
		jobs := st.List()
		json.NewEncoder(w).Encode(jobs)
	})

	fmt.Println("Scheduler running on :8080")
	http.ListenAndServe(":8080", nil)
}
