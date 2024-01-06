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
			log.Fatal(err)
		}

		// todo: handle other packets aside from PUBLISH
		go func() {
			switch RetrievePackageType(buf, n) {
			case PublishPacket:
				log.Println("publish packet")
			case SubscribePacket:
				log.Println("publish packet")
			default:
				log.Fatal("invalid packet type")
			}

			fmt.Printf("%x\n", buf[:n])
			pub := NewPublish(buf[:n])

			err = pub.Decode()
			if err != nil {
				log.Fatal(err)
			}

			resp, err := pub.Reply()
			if err != nil {
				log.Fatal(err)
			}

			// todo: send msg to subscribers
			// todo: save msg to history

			_, err = conn.Write(resp)
			if err != nil {
				log.Fatal(err)
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
		return nil, fmt.Errorf("failed to create connack")
	}

	return response, nil
}
