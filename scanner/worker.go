package scanner

import (
	"github.com/amine7536/quasar/utils"
)

// Worker executes a reverse lookup on a slice of ips
type Worker struct {
	ID            int
	Running       bool
	JobChannel    chan []string
	ResultChannel chan []string
	Done          chan int
	quit          chan bool
}

// NewWorker returns a new Worker
func NewWorker(workerID int, jobs chan []string, results chan []string, done chan int) Worker {
	return Worker{
		ID:            workerID,
		JobChannel:    jobs,
		ResultChannel: results,
		Done:          done,
		quit:          make(chan bool),
	}
}

// Start run the worker
func (w Worker) Start() {
	w.Running = true
	go func() {

		for {
			select {
			case IPs := <-w.JobChannel:
				// Send the return of fn in the ResultChannel
				for _, ip := range IPs {
					names, _ := utils.ResolveName(ip)
					// names := []string{"toto.com"}
					// time.Sleep(time.Millisecond * 1)
					w.ResultChannel <- append([]string{ip}, names...)
				}
				// Say that we have done the work
				w.Done <- w.ID

			case <-w.quit:
				// Stop working
				// log.Printf("Stopping WorkerID=%v", w.ID)
				w.Running = false
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
