package queue

import (
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	results := make(chan Job)
	defer close(results)

	tests := []struct {
		name       string
		maxWorkers int
	}{
		{
			name:       "create dispatcher with 1 worker",
			maxWorkers: 1,
		},
		{
			name:       "create dispatcher with 8 workers",
			maxWorkers: 8,
		},
		{
			name:       "create dispatcher with 100 workers",
			maxWorkers: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDispatcher(tt.maxWorkers, results)
			if d == nil {
				t.Fatal("NewDispatcher() returned nil")
			}
			if d.MaxWorkers != tt.maxWorkers {
				t.Errorf("NewDispatcher() MaxWorkers = %v, want %v", d.MaxWorkers, tt.maxWorkers)
			}
			if d.WorkerPool == nil {
				t.Error("NewDispatcher() WorkerPool is nil")
			}
			if d.JobQueue == nil {
				t.Error("NewDispatcher() JobQueue is nil")
			}
			if d.ResultQueue == nil {
				t.Error("NewDispatcher() ResultQueue is nil")
			}
		})
	}
}

func TestDispatcherRunStop(t *testing.T) {
	results := make(chan Job, 10)
	defer close(results)

	d := NewDispatcher(2, results)
	d.Run()

	// Give workers time to start
	time.Sleep(50 * time.Millisecond)

	// Verify workers were created
	if len(d.Workers) != 2 {
		t.Errorf("Run() created %d workers, want 2", len(d.Workers))
	}

	// Stop the dispatcher
	d.Stop()

	// Give time for cleanup
	time.Sleep(50 * time.Millisecond)
}

func TestDispatcherProcessJob(t *testing.T) {
	results := make(chan Job, 10)
	defer close(results)

	d := NewDispatcher(2, results)
	d.Run()
	defer d.Stop()

	// Give workers time to start
	time.Sleep(50 * time.Millisecond)

	// Send a test job
	testIP := "127.0.0.1"
	job := Job{IP: testIP}
	d.JobQueue <- job

	// Wait for result with timeout
	select {
	case result := <-results:
		if result.IP != testIP {
			t.Errorf("Result IP = %v, want %v", result.IP, testIP)
		}
		// Names may or may not be resolved depending on system configuration
		t.Logf("Job processed: IP=%s, Names=%v", result.IP, result.Names)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for job result")
	}
}

func TestDispatcherMultipleJobs(t *testing.T) {
	results := make(chan Job, 20)
	defer close(results)

	d := NewDispatcher(4, results)
	d.Run()
	defer d.Stop()

	// Give workers time to start
	time.Sleep(50 * time.Millisecond)

	// Send multiple jobs
	testJobs := []string{
		"127.0.0.1",
		"127.0.0.2",
		"127.0.0.3",
		"127.0.0.4",
		"127.0.0.5",
	}

	for _, ip := range testJobs {
		d.JobQueue <- Job{IP: ip}
	}

	// Collect results
	receivedJobs := 0
	timeout := time.After(10 * time.Second)

	for receivedJobs < len(testJobs) {
		select {
		case result := <-results:
			receivedJobs++
			found := false
			for _, ip := range testJobs {
				if result.IP == ip {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Received unexpected job with IP %v", result.IP)
			}
		case <-timeout:
			t.Fatalf("Timeout: received %d jobs, expected %d", receivedJobs, len(testJobs))
		}
	}

	if receivedJobs != len(testJobs) {
		t.Errorf("Received %d jobs, want %d", receivedJobs, len(testJobs))
	}
}

func TestNewWorker(t *testing.T) {
	workerPool := make(chan chan Job)
	results := make(chan Job)

	worker := NewWorker(1, workerPool, &results)

	if worker.ID != 1 {
		t.Errorf("NewWorker() ID = %v, want 1", worker.ID)
	}
	if worker.WorkerPool == nil {
		t.Error("NewWorker() WorkerPool is nil")
	}
	if worker.JobChannel == nil {
		t.Error("NewWorker() JobChannel is nil")
	}
	if worker.ResultChannel == nil {
		t.Error("NewWorker() ResultChannel is nil")
	}
}

func TestWorkerStartStop(t *testing.T) {
	workerPool := make(chan chan Job, 1)
	results := make(chan Job, 10)

	worker := NewWorker(1, workerPool, &results)
	worker.Start()

	// Give worker time to start
	time.Sleep(50 * time.Millisecond)

	// Worker should register itself in the pool
	select {
	case jobChan := <-workerPool:
		if jobChan == nil {
			t.Error("Worker registered nil job channel")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Worker did not register in pool")
	}

	// Stop the worker
	worker.Stop()
	time.Sleep(50 * time.Millisecond)
}

func TestWorkerProcessJob(t *testing.T) {
	workerPool := make(chan chan Job, 1)
	results := make(chan Job, 10)

	worker := NewWorker(1, workerPool, &results)
	worker.Start()
	defer worker.Stop()

	// Give worker time to start and register
	time.Sleep(50 * time.Millisecond)

	// Get the worker's job channel
	var jobChan chan Job
	select {
	case jobChan = <-workerPool:
	case <-time.After(1 * time.Second):
		t.Fatal("Worker did not register in pool")
	}

	// Send a job to the worker
	testIP := "127.0.0.1"
	jobChan <- Job{IP: testIP}

	// Wait for result
	select {
	case result := <-results:
		if result.IP != testIP {
			t.Errorf("Result IP = %v, want %v", result.IP, testIP)
		}
		t.Logf("Worker processed job: IP=%s, Names=%v", result.IP, result.Names)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for worker result")
	}
}

func TestJobStruct(t *testing.T) {
	job := Job{
		IP:    "192.168.1.1",
		Names: []string{"example.com", "test.com"},
	}

	if job.IP != "192.168.1.1" {
		t.Errorf("Job.IP = %v, want 192.168.1.1", job.IP)
	}
	if len(job.Names) != 2 {
		t.Errorf("Job.Names length = %v, want 2", len(job.Names))
	}
	if job.Names[0] != "example.com" {
		t.Errorf("Job.Names[0] = %v, want example.com", job.Names[0])
	}
}

// TestDispatcherConcurrency tests that multiple workers can process jobs concurrently
func TestDispatcherConcurrency(t *testing.T) {
	results := make(chan Job, 100)
	defer close(results)

	numWorkers := 10
	numJobs := 50

	d := NewDispatcher(numWorkers, results)
	d.Run()
	defer d.Stop()

	// Give workers time to start
	time.Sleep(100 * time.Millisecond)

	// Send jobs
	startTime := time.Now()
	for i := 0; i < numJobs; i++ {
		d.JobQueue <- Job{IP: "127.0.0.1"}
	}

	// Collect results
	receivedJobs := 0
	timeout := time.After(30 * time.Second)

	for receivedJobs < numJobs {
		select {
		case <-results:
			receivedJobs++
		case <-timeout:
			t.Fatalf("Timeout: received %d jobs, expected %d", receivedJobs, numJobs)
		}
	}

	duration := time.Since(startTime)
	t.Logf("Processed %d jobs with %d workers in %v", numJobs, numWorkers, duration)

	if receivedJobs != numJobs {
		t.Errorf("Received %d jobs, want %d", receivedJobs, numJobs)
	}
}
