package queue

import (
	"github.com/amine7536/reverse-scan/pkg/utils"
)

// Job represents a DNS lookup job
type Job struct {
	Names []string
	IP    string
}

// Worker executes a reverse lookup on a slice of ips
type Worker struct {
	WorkerPool    chan chan Job
	JobChannel    chan Job
	ResultChannel chan Job
	quit          chan bool
	ID            int
}

// NewWorker returns a new Worker
func NewWorker(id int, workerPool chan chan Job, resultQueue *chan Job) Worker {
	return Worker{
		ID:            id,
		WorkerPool:    workerPool,
		JobChannel:    make(chan Job),
		ResultChannel: *resultQueue,
		quit:          make(chan bool),
	}
}

// Start run the worker
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// Send the return of fn in the ResultChannel
				names, err := utils.ResolveName(job.IP)
				if err == nil {
					job.Names = names
				}
				w.ResultChannel <- job

			case <-w.quit:
				// Stop working
				// log.Printf("Stopping WorkerID=%v", w.ID)
				return

			}
		}
	}()
}

// Stop the Worker
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
