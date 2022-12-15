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

	// log.Println(client.GetInfo("39.144.103.149"))
	// log.Println(client.GetInfo("20.205.11.231"))
	// log.Println(client.GetInfo("222.64.171.253"))
	// log.Println(client.GetInfo("129.146.243.246"))
	// log.Println(client.GetInfo("192.168.189.125"))
	// log.Println(client.GetInfo("2600:4040:a912:a200:a438:9968:96d9:c3e4"))
	// log.Println(client.GetInfo("2600:387:1:809::3a"))
}
