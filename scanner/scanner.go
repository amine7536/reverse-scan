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

	reverseChan := make(chan []string)
	flushChan := make(chan bool)

	log.Printf("Resolving from %v to %v", config.StartIP, config.EndIP)
	log.Printf("Caluculated CIDR is %s", config.CIDR)

	mynet, _ := utils.GetHosts(config.CIDR)
	log.Printf("Number of IPs to scan: %v", len(mynet))

	file, err := os.Create(config.CSV)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	uiprogress.Start()                   // start rendering
	bar := uiprogress.AddBar(len(mynet)) // Add a new bar
	bar.AppendCompleted()
	bar.PrependElapsed()

	go rangeHosts(mynet, reverseChan, flushChan)

mainloop:
	for {
		select {
		case host := <-reverseChan:
			err := writer.Write(host)
			if err != nil {
				log.Fatal(err)
			}
			bar.Incr()
		case <-flushChan:
			writer.Flush()
			break mainloop
		}
	}

}

func rangeHosts(hosts []string, reverseChan chan []string, flushChan chan bool) {
	for _, ip := range hosts {
		// names, _ := utils.ResolveName(ip)
		// reverseChan <- append([]string{ip}, names...)
		names := []string{"toto.com"}
		reverseChan <- append([]string{ip}, names...)
		time.Sleep(time.Millisecond * 20)

		// err := writer.Write(append([]string{ip}, names...))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// log.Printf(ip, names)
	}

	flushChan <- true

}
