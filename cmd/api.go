package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/gussf/go-mqtt-server/domain/models"
)

var requestParser models.RequestParser

func startAPI(rp models.RequestParser, addr net.Addr) {
	requestParser = rp

	log.Printf("Starting server on %+v\n", addr)
	listener, err := net.Listen(addr.Network(), addr.String())
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

	resp, err := requestParser.ProcessConnectionRequest(buf)
	if err != nil {
		log.Fatal(err)
		return
	}

	n, err = conn.Write(resp)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer conn.Close()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("connection closed by peer")
				// todo remove from conn pool
				return
			}
			log.Print(err)
			return
		}

		go func() {
			resp, err := requestParser.ProcessRequest(buf, n)
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
