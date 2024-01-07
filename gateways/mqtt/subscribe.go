package mqtt

type Subscribe struct {
	raw            []byte
	VariableLength int
	Topic          string
}

func NewSubscribe(raw []byte) *Subscribe {
	return &Subscribe{
		raw: raw,
	}
}

func (s *Subscribe) Decode() error {
	return nil
}

func (s *Subscribe) Process() error {
	return nil
}

func (s *Subscribe) Reply() ([]byte, error) {
	fixedHeader := []byte{Suback, 0x03, 0x00, 0x00, 0x00}
	resp := fixedHeader

	return resp, nil
}
