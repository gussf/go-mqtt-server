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

	resp, err := parseMQTT(buf)
	if err != nil {
		log.Fatal(err)
	}

	n, err = conn.Write(resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n, resp)
}

func parseMQTT(packet []byte) ([]byte, error) {
	if len(packet) == 0 {
		return nil, errors.New("mqtt packet is empty")
	}

	var msg MQTT
	if packet[0] >= 0x10 && packet[0] < 0x20 {
		log.Println("Conn packet")
		msg = NewConnect(packet)
	} else {
		log.Println("no")
	}

	err := msg.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to decode packet: %w", err)
	}

	response, err := msg.Reply()
	if err != nil {
		return nil, fmt.Errorf("failed to create message reply")
	}
	// 10 2d 00 04 4d 51 54 54 04 02 00 3c0021706f73746d616e2d6d7174742d636c69656e742d31373033393732353034373536
	return response, nil
}
