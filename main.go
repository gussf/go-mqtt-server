package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp4", ":8001")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("failed to accept connection", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("failed to read ", err)
	}
	fmt.Printf("%x\n", buf[:n])

	resp, err := parseConnectionPacket(buf)
	if err != nil {
		log.Fatal(err)
	}

	n, err = conn.Write(resp)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		go func() {
			resp, err := processRequest(buf, n)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = conn.Write(resp)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func parseConnectionPacket(packet []byte) ([]byte, error) {
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

func processRequest(buf []byte, size int) ([]byte, error) {
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
