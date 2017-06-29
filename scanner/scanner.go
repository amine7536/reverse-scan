package scanner

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/gosuri/uiprogress"

	"bitbucket.org/aminebenseddik/reverse-scan/conf"
	"bitbucket.org/aminebenseddik/reverse-scan/utils"
)

func Start(config *conf.Config) {

	hosts, _ := utils.GetHosts(config.CIDR)
	jobsChan := make(chan []string)
	resultsChan := make(chan []string)
	doneChan := make(chan int, config.WORKERS)

	// var wg sync.WaitGroup

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
	bar := uiprogress.AddBar(len(hosts))
	bar.AppendCompleted()
	bar.PrependElapsed()

	// Split hosts list
	var workers []Worker
	var stoppedCount = 0

	for a, b := range utils.SplitSlice(hosts, config.WORKERS) {
		workers = append(workers, NewWorker(a, jobsChan, resultsChan, doneChan))
		// log.("Starting WorkerID=%v with slice lenght=%v", a, len(b))
		workers[a].Start()
		workers[a].JobChannel <- b
	}

mainloop:
	for {
		select {
		case res := <-resultsChan:
			err := writer.Write(append(res))
			if err != nil {
				log.Fatal(err)
			}
			writer.Flush()
			bar.Incr()
		case id := <-doneChan:
			writer.Flush()
			workers[id].Stop()

			stoppedCount++
			if stoppedCount == config.WORKERS {
				time.Sleep(time.Second) // wait for a second for all the go routines to finish
				uiprogress.Stop()
				break mainloop
			}

		}
	}
}
