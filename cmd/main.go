package main

import (
	"log"
	"net"

	"github.com/gussf/go-mqtt-server/gateways/mqtt"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp4", "localhost:8001")
	if err != nil {
		log.Fatal(err)
	}

	rp := mqtt.NewMQTTParser()

	api := NewAPI(rp, addr)
	api.Start()
}
