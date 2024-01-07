package main

import (
	"log"
	"net"

	"github.com/gussf/go-mqtt-server/cmd/api"
)

type Server struct {
	API   *api.API
	Addr net.Addr
}

func NewServer(api *api.API, addr net.Addr) *Server {
	return &Server{
		API:   api,
		Addr: addr,
	}
}

func (s *Server) Start() {
	s.API.Start(s.Addr)
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp4", "localhost:8001")
	if err != nil {
		log.Fatal(err)
	}

	s := InitializeServer(addr)
	s.Start()
}
