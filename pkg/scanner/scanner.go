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

	hosts, _ := utils.GetHosts(c.CIDR)
	results := make(chan queue.Job)
	defer close(results)

	log.Printf("Resolving from %v to %v", c.StartIP, c.EndIP)
	log.Printf("Caluculated CIDR is %s", c.CIDR)
	log.Printf("Number of IPs to scan: %v", len(hosts))
	log.Printf("Starting %v Workers", c.WORKERS)

	file, err := os.Create(c.CSV)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
	r := 0
resultLoop:
	for {
		select {
		case job := <-results:
			err := writer.Write(append([]string{job.IP}, job.Names...))
			if err != nil {
				log.Fatal(err)
			}
			writer.Flush()
			bar.Incr()
			r++
			if r == len(hosts) {
				// We got all jobs back, we stop
				break resultLoop
			}
		}
	}
}
