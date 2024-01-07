package mqtt

type Connect struct {
	raw []byte
}

func NewConnect(raw []byte) *Connect {
	return &Connect{
		raw: raw,
	}
}

func (c *Connect) Decode() error {
	return nil
}

func (c *Connect) Process() error {
	return nil
}

func (c *Connect) Reply() ([]byte, error) {
	resp := []byte{Connack}
	resp = append(resp, ConnackLength, 0x00, 0x00, 0x00)

	return resp, nil
}
