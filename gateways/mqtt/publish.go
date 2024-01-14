package mqtt

import (
	"errors"
	"log"

	"github.com/gussf/go-mqtt-server/domain/models"
	"github.com/gussf/go-mqtt-server/domain/usecases/publisher"
)

type PublishHandler struct {
	uc             publisher.Usecase
	raw            []byte
	VariableLength int
	Topic          []byte
	Payload        []byte
}

func NewPublishHandler(uc publisher.Usecase) *PublishHandler {
	return &PublishHandler{
		uc: uc,
	}
}

func (p *PublishHandler) Decode(buf []byte) error {
	if len(buf) < 2 {
		return errors.New("invalid publish length")
	}
	p.raw = buf

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
	p.Topic = topicBytes

	// handle payload
	remaining := p.VariableLength - (int(msbLength) + int(lsbLength) + 2)
	payloadBytes := make([]byte, remaining)
	for i := 0; i < remaining; i++ {
		payloadBytes[i] = p.raw[currPos]
		currPos++
	}
	p.Payload = payloadBytes

	log.Printf("Publish:%#v\n", *p)
	return nil
}

func (p *PublishHandler) Encode() []byte {
	encoded := make([]byte, 4) // hardcoded bc im lazy
	encoded[0] = PublishRequest

	variableLen := len(p.Payload) + len(p.Topic)
	encoded[1] = byte(variableLen)
	encoded[2] = 0x00
	encoded[3] = byte(len(p.Topic))
	encoded = append(encoded, p.Topic...)
	encoded = append(encoded, p.Payload...)

	return encoded
}

// Process posts the published message to all subscribers
func (p *PublishHandler) Process(conn models.Connection) error {
	sub := models.Subscription{
		Topic: string(p.Topic),
		Conn:  conn,
	}

	payload := p.Encode()
	err := p.uc.PublishToSubscribers(sub, payload)
	if err != nil {
		return err
	}

	return nil
}

func (p *PublishHandler) Reply(models.Connection) ([]byte, error) {
	fixedHeader := []byte{Puback, 0x02, 0x00, 0x00}

	resp := fixedHeader

	return resp, nil
}
