package models

import "net"

type RequestParser interface {
	ProcessRequest(buf []byte, size int, conn net.Conn) ([]byte, error)
	ProcessConnectionRequest(packet []byte) ([]byte, error)
}
