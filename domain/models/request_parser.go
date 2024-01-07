package models

type RequestParser interface {
	ProcessRequest(buf []byte, size int) ([]byte, error)
	ProcessConnectionRequest(packet []byte) ([]byte, error)
}
