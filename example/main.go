package main

import (
	"log"
	"time"

	"github.com/coreservice-io/geo_ip/lib"
)

func main() {

	client, err := lib.NewClient("./example", true, func(logstr string) {
		log.Println(logstr)
	})

	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(client.GetInfo("176.119.148.39"))
	log.Println(client.GetInfo("103.140.9.84"))
	log.Println(client.GetInfo("94.72.140.55"))
	log.Println(client.GetInfo("103.177.80.154"))
	log.Println(client.GetInfo("154.19.185.171"))
	log.Println(client.GetInfo("202.81.232.120"))
	log.Println(client.GetInfo("38.6.229.3"))

	log.Println(client.GetInfo("123.118.103.129"))
	log.Println(client.GetInfo("185.100.232.166"))
	log.Println(client.GetInfo("107.164.105.2"))
	log.Println(client.GetInfo("107.164.105.30"))

	time.Sleep(5 * time.Hour)

}
