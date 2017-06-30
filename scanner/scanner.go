package scanner

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gosuri/uiprogress"

	"bitbucket.org/aminebenseddik/reverse-scan/conf"
	"bitbucket.org/aminebenseddik/reverse-scan/utils"
)

func Start(config *conf.Config) {

	hosts, _ := utils.GetHosts(config.CIDR)
	results := make(chan Job)
	defer close(results)

	log.Printf("Resolving from %v to %v", config.StartIP, config.EndIP)
	log.Printf("Caluculated CIDR is %s", config.CIDR)
	log.Printf("Number of IPs to scan: %v", len(hosts))
	log.Printf("Starting %v Workers", config.WORKERS)

	file, err := os.Create(config.CSV)
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

	dispatch := NewDispatcher(config.WORKERS, results)
	dispatch.Run()
	defer dispatch.Stop()
	// time.Sleep(time.Second * 10)

	// Send Jobs to Dispatch
	for _, ip := range hosts {
		work := Job{IP: ip}
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
