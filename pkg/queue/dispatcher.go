package queue

// Dispatcher dispatchs jobs to worker
type Dispatcher struct {
	MaxWorkers  int
	WorkerPool  chan chan Job
	JobQueue    chan Job
	ResultQueue chan Job
	quit        chan bool
	Workers     []Worker
}

// NewDispatcher returns a new dispatcher
func NewDispatcher(maxWorkers int, results chan Job) *Dispatcher {

	return &Dispatcher{
		MaxWorkers:  maxWorkers,
		WorkerPool:  make(chan chan Job, maxWorkers),
		JobQueue:    make(chan Job),
		ResultQueue: results,
		quit:        make(chan bool),
	}
}

// Run starts the dispatcher
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i, d.WorkerPool, &d.ResultQueue)
		d.Workers = append(d.Workers, worker)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)

		case <-d.quit:
			for _, w := range d.Workers {
				w.Stop()
			}
			return
		}
	}
}

// Stop stops the dispatcher and all workers
func (d *Dispatcher) Stop() {
	go func() {
		d.quit <- true
	}()
}
