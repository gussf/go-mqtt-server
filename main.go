package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gussf/go-mqtt-server/domain/models"
	"github.com/gussf/go-mqtt-server/gateways/mqtt"
)

var rp models.RequestParser

func main() {
	rp = mqtt.NewMQTTParser()
	listener, err := net.Listen("tcp4", ":8001")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("failed to accept connection", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("failed to read ", err)
	}
	fmt.Printf("%x\n", buf[:n])

	resp, err := rp.ProcessConnectionRequest(buf)
	if err != nil {
		log.Fatal(err)
	}

	n, err = conn.Write(resp)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		go func() {
			resp, err := rp.ProcessRequest(buf, n)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = conn.Write(resp)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}
