package mqtt

import (
	"errors"
	"fmt"
	"log"

	"github.com/gussf/go-mqtt-server/domain/models"
)

type MQTTParser struct {
	pool models.SubscriptionPool
}

func NewMQTTParser() *MQTTParser {
	return &MQTTParser{
		pool: models.NewSubscriptionPool()}
}

func (m *MQTTParser) ProcessConnectionRequest(packet []byte) ([]byte, error) {
	if len(packet) == 0 {
		return nil, errors.New("connection packet is empty")
	}
	packetType := RetrievePackageType(packet, len(packet))
	if packetType != ConnectionPacket {
		return nil, fmt.Errorf("not a connection packet")
	}

	msg := NewConnect(packet)
	err := msg.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to decode packet: %w", err)
	}

	response, err := msg.Reply()
	if err != nil {
		return nil, fmt.Errorf("failed to create connack: %w", err)
	}

	return response, nil
}

func (m *MQTTParser) ProcessRequest(buf []byte, size int) ([]byte, error) {
	fmt.Printf("Received: %x\n", buf[:size])
	var resp []byte
	var msg MQTT
	switch RetrievePackageType(buf, size) {
	case PublishPacket:
		log.Println("publish packet")
		msg = NewPublish(buf[:size])
	case SubscribePacket:
		log.Println("subscribe packet")
		msg = NewSubscribe(buf[:size])
	case PingPacket:
		log.Println("ping packet")
		msg = NewPing(buf[:size])
	default:
		log.Fatal("invalid packet type")
	}

	err := msg.Decode()
	if err != nil {
		return nil, err
	}

	err = msg.Process()
	if err != nil {
		return nil, err
	}

	resp, err = msg.Reply()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Sent: %x\n", resp)

	return resp, nil
}
