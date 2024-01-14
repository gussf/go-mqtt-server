package mqtt

import "github.com/gussf/go-mqtt-server/domain/models"

type PingHandler struct {
	raw []byte
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (p *PingHandler) Decode(buf []byte) error {
	return nil
}

func (p *PingHandler) Process(conn models.Connection) error {
	return nil
}

func (p *PingHandler) Reply(conn models.Connection) ([]byte, error) {
	fixedHeader := []byte{PingResp, 0x00}
	resp := fixedHeader

	return resp, nil
}
