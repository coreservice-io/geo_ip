package main

import (
	"log"

	"github.com/coreservice-io/geo_ip/local"
)

func main() {

	local, err := local.NewIp2L("./example/geo_ip.csv")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(local.GetLocalInfo("129.146.243.246"))
	log.Println(local.GetLocalInfo("192.168.189.125"))
}
