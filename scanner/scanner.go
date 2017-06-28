package scanner

import (
	"encoding/csv"
	"log"
	"os"

	"bitbucket.org/aminebenseddik/reverse-scan/conf"
	"bitbucket.org/aminebenseddik/reverse-scan/utils"
)

func Start(config *conf.Config) {

	reverseChan := make(chan []string)

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

	go rangeHosts(mynet, reverseChan)

	for {
		select {
		case host := <-reverseChan:
			err := writer.Write(host)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func rangeHosts(hosts []string, reverseChan chan []string) {
	for _, ip := range hosts {
		names, _ := utils.ResolveName(ip)
		reverseChan <- append([]string{ip}, names...)

		// err := writer.Write(append([]string{ip}, names...))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// log.Printf(ip, names)
	}

}
