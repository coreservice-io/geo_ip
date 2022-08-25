package main

import (
	"log"

	"github.com/coreservice-io/geo_ip/lib"
)

func main() {

	client, err := lib.NewClient("./example/geo_ip.csv")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(client.GetInfo("129.146.243.246"))
	log.Println(client.GetInfo("192.168.189.125"))
}
