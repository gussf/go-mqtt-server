package main

type MQTT interface {
	Decode() error
	Reply() ([]byte, error)
}

const (
	ConnackSuccess = 0x20
	ConnackLength  = 0x03
)

type Connect struct {
	raw          []byte
	ProtocolName string
}

func NewConnect(raw []byte) *Connect {
	return &Connect{
		raw: raw,
	}
}

func (c *Connect) Decode() error {
	return nil
}

func (c *Connect) Reply() ([]byte, error) {
	resp := []byte{ConnackSuccess}
	resp = append(resp, ConnackLength, 0x00, 0x00, 0x00)

	return resp, nil
}
