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
	Topic          string
	Payload        string
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
func (p *PublishHandler) Process(conn models.Connection) error {
	sub := models.Subscription{
		Topic: p.Topic,
		Conn:  conn,
	}

	err := p.uc.PublishToSubscribers(sub, p.Payload)
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
