// Package scanner performs reverse DNS lookups on IP ranges
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

	log.Printf("Resolving from %v to %v", c.StartIP, c.EndIP)
	log.Printf("Calculated CIDR is %s", c.CIDR)
	log.Printf("Number of IPs to scan: %v", len(hosts))
	log.Printf("Starting %v Workers", c.WORKERS)

	file, err := os.Create(c.CSV)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}

	writer := csv.NewWriter(file)

	uiprogress.Start()
	bar := uiprogress.AddBar(len(hosts))
	bar.AppendCompleted()
	bar.PrependElapsed()

	dispatch := queue.NewDispatcher(c.WORKERS, results)
	dispatch.Run()

	// Send Jobs to Dispatch
	for _, ip := range hosts {
		work := queue.Job{IP: ip}
		dispatch.JobQueue <- work
	}

	// Wait for results
	for r := 0; r < len(hosts); r++ {
		job := <-results
		if err := writer.Write(append([]string{job.IP}, job.Names...)); err != nil {
			writer.Flush()
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("Warning: failed to close file: %v", closeErr)
			}
			uiprogress.Stop()
			dispatch.Stop()
			log.Fatalf("Failed to write result: %v", err)
		}
		writer.Flush()
		bar.Incr()
	}

	writer.Flush()
	if err := file.Close(); err != nil {
		log.Printf("Warning: failed to close file: %v", err)
	}
	uiprogress.Stop()
	dispatch.Stop()
}
