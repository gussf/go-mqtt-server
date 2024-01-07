package mqtt

type Ping struct {
	raw []byte
}

func NewPing(raw []byte) *Ping {
	return &Ping{
		raw: raw,
	}
}

func (p *Ping) Decode() error {
	return nil
}

func (p *Ping) Process() error {
	return nil
}

func (p *Ping) Reply() ([]byte, error) {
	fixedHeader := []byte{PingResp, 0x00}
	resp := fixedHeader

	return resp, nil
}
