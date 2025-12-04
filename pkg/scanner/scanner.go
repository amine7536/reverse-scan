package scanner

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gosuri/uiprogress"

	"github.com/amine7536/reverse-scan/pkg/config"
	"github.com/amine7536/reverse-scan/pkg/queue"
	"github.com/amine7536/reverse-scan/pkg/utils"
)

// Start scanner
func Start(c *config.Config) {
	hosts, err := utils.GetHosts(c.CIDR)
	if err != nil {
		log.Fatalf("Failed to get hosts: %v", err)
	}

	results := make(chan queue.Job)
	defer close(results)

	log.Printf("Resolving from %v to %v", c.StartIP, c.EndIP)
	log.Printf("Calculated CIDR is %s", c.CIDR)
	log.Printf("Number of IPs to scan: %v", len(hosts))
	log.Printf("Starting %v Workers", c.WORKERS)

	file, err := os.Create(c.CSV)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: failed to close file: %v", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	uiprogress.Start()
	defer uiprogress.Stop()
	bar := uiprogress.AddBar(len(hosts))
	bar.AppendCompleted()
	bar.PrependElapsed()

	dispatch := queue.NewDispatcher(c.WORKERS, results)
	dispatch.Run()
	defer dispatch.Stop()
	// time.Sleep(time.Second * 10)

	// Send Jobs to Dispatch
	for _, ip := range hosts {
		work := queue.Job{IP: ip}
		dispatch.JobQueue <- work
	}

	// Wait for results
	for r := 0; r < len(hosts); r++ {
		job := <-results
		err := writer.Write(append([]string{job.IP}, job.Names...))
		if err != nil {
			log.Fatalf("Failed to write result: %v", err)
		}
		writer.Flush()
		bar.Incr()
	}
}
