package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/gussf/go-mqtt-server/domain/models"
)

type API struct {
	addr net.Addr
	rp   models.RequestParser
}

func NewAPI(rp models.RequestParser, addr net.Addr) *API {
	return &API{
		addr: addr,
		rp: rp,
	}
}

func (a *API) Start() {
	log.Printf("Starting server on %+v\n", a.addr)

	listener, err := net.Listen(a.addr.Network(), a.addr.String())
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
		go a.handleConnection(conn)
	}
}

func (a *API) handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("failed to read ", err)
	}
	fmt.Printf("%x\n", buf[:n])

	resp, err := a.rp.ProcessConnectionRequest(buf)
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
			resp, err := a.rp.ProcessRequest(buf, n)
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
