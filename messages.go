package main

// todo 
// message package

import (
	"errors"
	"log"
)

type MQTT interface {
	Decode() error
	Process() error
	Reply() ([]byte, error)
}

type PacketType int

const (
	ConnectionPacket PacketType = iota
	PublishPacket 
	SubscribePacket 
	PingPacket
	UnknownPacket 
)

const (
	ConnectionRequest = 0x10
	Connack = 0x20
	ConnackLength  = 0x03
	PublishRequest = 0x30
	Puback  = 0x40
	SubscribeRequest = 0x80
	Suback= 0x90
	PingReq = 0xC0
	PingResp = 0xD0
)

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

func (p *Publish) Process() error {
	return nil
}


func (p *Publish) Reply() ([]byte, error) {
	fixedHeader := []byte{Puback, 0x02, 0x00, 0x00}

	resp := fixedHeader

	return resp, nil
}

type Subscribe struct {
	raw            []byte
	VariableLength int
	Topic          string
}

func NewSubscribe(raw []byte) *Subscribe{
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

type Ping struct {
	raw            []byte
}

func NewPing(raw []byte) *Ping{
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

func RetrievePackageType(buf []byte, size int) PacketType {
	if size <= 1 { 
		return UnknownPacket
	}

	req := buf[0]
	if req >= ConnectionRequest && req < Connack {
		return ConnectionPacket
	} else if req >= PublishRequest && req < Puback {
		return PublishPacket
	} else if req >= SubscribeRequest && req < Suback {
		return SubscribePacket
	} else if req >= PingReq && req < PingResp {
		return PingPacket
	}

	return UnknownPacket
}
