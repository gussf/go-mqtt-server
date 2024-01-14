package mqtt

import (
	"errors"
	"log"
	"net"

	"github.com/gussf/go-mqtt-server/domain/models"
	"github.com/gussf/go-mqtt-server/domain/usecases/publisher"
)

type SubscribeHandler struct {
	raw            []byte
	uc             publisher.Usecase
	conn           net.Conn
	VariableLength int
	Topic          string
	TopicId        []byte
}

func NewSubscribeHandler(uc publisher.Usecase) *SubscribeHandler {
	return &SubscribeHandler{
		uc:      uc,
		TopicId: []byte{},
	}
}

func (s *SubscribeHandler) Decode(buf []byte) error {
	if len(buf) < 2 {
		return errors.New("invalid subscribe length")
	}
	s.raw = buf

	SubscribeLengthPos := 1
	s.VariableLength = int(s.raw[SubscribeLengthPos])

	if s.VariableLength <= 0 {
		return errors.New("invalid subscribe variable header length")
	}
	currPos := SubscribeLengthPos
	currPos++

	SubscribeTopicIdLength := 2
	topicIdBytes := make([]byte, SubscribeTopicIdLength)
	for i := 0; i < int(SubscribeTopicIdLength); i++ {
		topicIdBytes[i] = s.raw[currPos]
		currPos++
	}
	s.TopicId = topicIdBytes

	// handle MSB
	msbLength := int(s.raw[currPos])
	currPos++
	for i := 0; i < int(msbLength); i++ {
		// idk
		currPos++
	}

	// handle LSB
	lsbLength := int(s.raw[currPos])
	currPos++

	topicBytes := make([]byte, lsbLength)
	for i := 0; i < int(lsbLength); i++ {
		topicBytes[i] = s.raw[currPos]
		currPos++
	}
	s.Topic = string(topicBytes)
	log.Printf("Subscribe:%#v\n", *s)

	return nil
}

func (s *SubscribeHandler) Process(conn models.Connection) error {
	sub := models.Subscription{
		Topic: s.Topic,
		Conn:  conn,
	}

	err := s.uc.AddConnectionToTopicPool(sub)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscribeHandler) Reply(conn models.Connection) ([]byte, error) {
	fixedHeader := []byte{Suback, 0x03, 0x00, 0x00, 0x00}
	resp := fixedHeader

	return resp, nil
}
