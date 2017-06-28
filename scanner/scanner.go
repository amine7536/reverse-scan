package scanner

import (
	"encoding/csv"
	"log"
	"os"

	"bitbucket.org/aminebenseddik/reverse-scan/conf"
	"bitbucket.org/aminebenseddik/reverse-scan/utils"
)

func Start(config *conf.Config) {

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

	for _, ip := range mynet {
		names, _ := utils.ResolveName(ip)

		// err := writer.Write(append([]string{ip}, names...))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		log.Printf(ip, names)
	}

}
