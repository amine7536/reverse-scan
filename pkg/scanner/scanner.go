package scanner

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/amine7536/reverse-scan/pkg/config"
	"github.com/amine7536/reverse-scan/pkg/queue"
	"github.com/amine7536/reverse-scan/pkg/utils"
	"github.com/gosuri/uiprogress"
)

// Start scanner
func Start(c *config.Config) error {
	var err error
	hosts, _ := utils.GetHosts(c.CIDR)
	results := make(chan queue.Job)
	defer close(results)

	log.Printf("Resolving from %v to %v", c.StartIP, c.EndIP)
	log.Printf("Calculated CIDR is %s", c.CIDR)
	log.Printf("Number of IPs to scan: %v", len(hosts))
	log.Printf("Starting %v Workers", c.WORKERS)

	file, err := os.Create(c.CSV)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
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

	// Send Jobs to Dispatch
	for _, ip := range hosts {
		work := queue.Job{IP: ip}
		dispatch.JobQueue <- work
	}

	// Wait for results
	r := 0
	for job := range results {
		err := writer.Write(append([]string{job.IP}, job.Names...))
		if err != nil {
			log.Println(err)
			return err
		}
		writer.Flush()
		bar.Incr()
		r++
		if r == len(hosts) {
			// We got all jobs back, we stop
			break
		}
	}

	return nil
}
