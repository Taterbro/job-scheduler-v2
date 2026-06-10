package main

import (
	"fmt"
	"net/http"

	"github.com/taterbro/job-scheduler-v2/internal/handler"
	"github.com/taterbro/job-scheduler-v2/internal/routes"
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

	http.HandleFunc("/jobs", routes.CreateJob(sched))
	http.HandleFunc("/jobs/list", routes.ListJobs(st))

	fmt.Println("Scheduler running on :8080")
	http.ListenAndServe(":8080", nil)
}
