package main

import (
	"log"

	"github.com/coreservice-io/geo_ip/lib"
)

func main() {

	client, err := lib.NewClient(
		"./example/country_ipv4.csv",
		"./example/country_ipv6.csv",
		"./example/isp_ipv4.csv",
		"./example/isp_ipv6.csv")

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(client.GetInfo("185.7.108.0"))
	log.Println(client.GetInfo("39.144.103.149"))
	log.Println(client.GetInfo("20.205.11.231"))
	log.Println(client.GetInfo("222.64.171.253"))
	log.Println(client.GetInfo("129.146.243.246"))
	log.Println(client.GetInfo("192.168.189.125"))
	log.Println(client.GetInfo("2600:4040:a912:a200:a438:9968:96d9:c3e4"))
	log.Println(client.GetInfo("2600:387:1:809::3a"))

}
