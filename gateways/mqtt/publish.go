package mqtt

import (
	"errors"
	"log"
)

type Publish struct {
	raw            []byte
	VariableLength int
	Topic          string
	Payload        string
}

func NewPublish(raw []byte) *Publish {
	return &Publish{
		raw: raw,
	}
}

func (p *Publish) Decode() error {
	if len(p.raw) < 2 {
		return errors.New("invalid publish length")
	}

	PublishLengthPos := 1
	p.VariableLength = int(p.raw[PublishLengthPos])

	if p.VariableLength <= 0 {
		return errors.New("invalid publish variable header length")
	}
	currPos := PublishLengthPos
	currPos++

	// todo handle other flags such as QoS

	// handle MSB
	msbLength := int(p.raw[currPos])
	currPos++
	for i := 0; i < int(msbLength); i++ {
		// idk
		currPos++
	}

	// handle LSB
	lsbLength := int(p.raw[currPos])
	currPos++

	topicBytes := make([]byte, lsbLength)
	for i := 0; i < int(lsbLength); i++ {
		topicBytes[i] = p.raw[currPos]
		currPos++
	}
	p.Topic = string(topicBytes)

	// handle payload
	remaining := p.VariableLength - (int(msbLength) + int(lsbLength) + 2)
	payloadBytes := make([]byte, remaining)
	for i := 0; i < remaining; i++ {
		payloadBytes[i] = p.raw[currPos]
		currPos++
	}
	p.Payload = string(payloadBytes)

	log.Printf("Publish:%#v\n", *p)
	return nil
}

// Process posts the published message to all subscribers
func (p *Publish) Process() error {
	return nil
}

func (p *Publish) Reply() ([]byte, error) {
	fixedHeader := []byte{Puback, 0x02, 0x00, 0x00}

	resp := fixedHeader

	return resp, nil
}
