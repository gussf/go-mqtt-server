//go:build wireinject

package main

import (
	"net"

	"github.com/google/wire"
	"github.com/gussf/go-mqtt-server/cmd/api"
	"github.com/gussf/go-mqtt-server/gateways/mqtt"
)

func InitializeServer(addr net.Addr) *Server {
	wire.Build(mqtt.NewMQTTParser, api.NewAPI, NewServer)
	return &Server{}
}
