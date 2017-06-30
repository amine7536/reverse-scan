package scanner

import (
	"bitbucket.org/aminebenseddik/reverse-scan/utils"
)

type Job struct {
	IP    string
	Names []string
}

// Worker executes a reverse lookup on a slice of ips
type Worker struct {
	ID            int
	WorkerPool    chan chan Job
	JobChannel    chan Job
	ResultChannel chan Job
	quit          chan bool
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
				job.Names, _ = utils.ResolveName(job.IP)
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
