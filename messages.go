package main

import (
	"errors"
	"log"
)

type MQTT interface {
	Decode() error
	Reply() ([]byte, error)
}

const (
	ConnackSuccess = 0x20
	PubackSuccess  = 0x40
	ConnackLength  = 0x03
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

func (c *Connect) Reply() ([]byte, error) {
	resp := []byte{ConnackSuccess}
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

func (p *Publish) Reply() ([]byte, error) {
	fixedHeader := []byte{PubackSuccess, 0x02, 0x00, 0x00}

	resp := fixedHeader

	return resp, nil
}
