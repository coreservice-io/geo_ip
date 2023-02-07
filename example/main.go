package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coreservice-io/geo_ip/lib"
)

const geo_ip_update_key = ""

func main() {

	update_err := lib.Update(geo_ip_update_key, "./example", func(s string) {
		fmt.Println("logstr:", s)
	})
	if update_err != nil {
		panic(update_err)
	}

	//////////
	client, err := lib.NewClient(geo_ip_update_key, "0.0.24", "./example", func(log_str string) {
		fmt.Println("log_str:" + log_str)
	}, func(err_log_str string) {
		fmt.Println("err_log_str:" + err_log_str)
	})

	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(client.GetInfo("104.233.16.169"))
	log.Println(client.GetInfo("5.78.52.174"))
	log.Println(client.GetInfo("116.227.21.107"))

	time.Sleep(30 * time.Second)

	log.Println(client.GetInfo("172.104.160.0"))

	time.Sleep(30 * time.Hour)
}
