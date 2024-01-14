package mqtt

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/gussf/go-mqtt-server/domain/models"
	"github.com/gussf/go-mqtt-server/domain/usecases/publisher"
)

type MQTTParser struct {
	uc               publisher.Usecase
	publishHandler   *PublishHandler
	subscribeHandler *SubscribeHandler
	pingHandler      *PingHandler
}

func NewMQTTParser(uc publisher.Usecase) models.RequestParser {
	return &MQTTParser{
		uc:               uc,
		publishHandler:   NewPublishHandler(uc),
		subscribeHandler: NewSubscribeHandler(),
		pingHandler:      NewPingHandler(),
	}
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

func (m *MQTTParser) ProcessRequest(buf []byte, size int, conn net.Conn) ([]byte, error) {
	if size == 0 {
		return nil, errors.New("invalid buffer")
	}
	fmt.Printf("Received: %x\n", buf[:size])

	domainConn := models.Connection{Conn: conn}
	var resp []byte
	var msg MQTT
	switch RetrievePackageType(buf, size) {
	case PublishPacket:
		log.Println("publish packet")
		msg = m.publishHandler
	case SubscribePacket:
		log.Println("subscribe packet")
		msg = m.subscribeHandler
	case PingPacket:
		log.Println("ping packet")
		msg = m.pingHandler
	default:
		log.Fatal("invalid packet type")
	}

	err := msg.Decode(buf)
	if err != nil {
		return nil, err
	}

	err = msg.Process(domainConn)
	if err != nil {
		return nil, err
	}

	resp, err = msg.Reply(domainConn)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Sent: %x\n", resp)

	return resp, nil
}
