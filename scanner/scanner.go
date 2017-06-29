package scanner

import (
	"encoding/csv"
	"log"
	"os"
	"sync"

	"github.com/gosuri/uiprogress"

	"bitbucket.org/aminebenseddik/reverse-scan/conf"
	"bitbucket.org/aminebenseddik/reverse-scan/utils"
)

const (
	WORKERS = 4
)

func Start(config *conf.Config) {

	hosts, _ := utils.GetHosts(config.CIDR)
	jobsChan := make(chan []string)
	resultsChan := make(chan []string)
	doneChan := make(chan int, WORKERS)

	var wg sync.WaitGroup

	log.Printf("Resolving from %v to %v", config.StartIP, config.EndIP)
	log.Printf("Caluculated CIDR is %s", config.CIDR)
	log.Printf("Number of IPs to scan: %v", len(hosts))

	file, err := os.Create(config.CSV)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	//defer writer.Flush()

	uiprogress.Start()
	bar := uiprogress.AddBar(len(hosts))
	bar.AppendCompleted()
	bar.PrependElapsed()

	// Split hosts list
	var workers []Worker

	for a, b := range utils.SplitSlice(hosts, WORKERS) {
		wg.Add(1)
		workers = append(workers, NewWorker(a, jobsChan, resultsChan, doneChan))
		workers[a].Start(wg)
		workers[a].JobChannel <- b
	}

	var stoppedCount = 0

	// mainloop:
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
			// if stoppedCount == WORKERS {
			// 	break mainloop
			// }

		}
	}

	wg.Wait()
	// break mainloop

}
