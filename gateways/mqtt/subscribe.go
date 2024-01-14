package mqtt

import (
	"net"

	"github.com/gussf/go-mqtt-server/domain/models"
)

type SubscribeHandler struct {
	raw            []byte
	conn           net.Conn
	VariableLength int
	Topic          string
}

func NewSubscribeHandler() *SubscribeHandler {
	return &SubscribeHandler{}
}

func (s *SubscribeHandler) Decode(buf []byte) error {
	return nil
}

func (s *SubscribeHandler) Process(conn models.Connection) error {
	return nil
}

func (s *SubscribeHandler) Reply(conn models.Connection) ([]byte, error) {
	fixedHeader := []byte{Suback, 0x03, 0x00, 0x00, 0x00}
	resp := fixedHeader

	return resp, nil
}
