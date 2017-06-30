package scanner

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	MaxWorkers  int
	WorkerPool  chan chan Job
	JobQueue    chan Job
	ResultQueue chan Job
	quit        chan bool
	Workers     []Worker
}

func NewDispatcher(maxWorkers int, results chan Job) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	jobs := make(chan Job)

	return &Dispatcher{
		MaxWorkers:  maxWorkers,
		WorkerPool:  pool,
		JobQueue:    jobs,
		ResultQueue: results,
	}
}

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

func (d *Dispatcher) Stop() {
	go func() {
		d.quit <- true
	}()
}
