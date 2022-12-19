package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coreservice-io/geo_ip/lib"
)

func main() {

	client, err := lib.NewClient("0.0.10", "./example", func(log_str string) {
		fmt.Println("log_str:" + log_str)
	}, func(err_log_str string) {
		fmt.Println("err_log_str:" + err_log_str)
	})

	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(client.GetInfo("172.104.160.0"))
	log.Println(client.GetInfo("178.239.197.0"))

	time.Sleep(30 * time.Second)

	log.Println(client.GetInfo("172.104.160.0"))
	log.Println(client.GetInfo("178.239.197.0"))

	time.Sleep(30 * time.Hour)
}
