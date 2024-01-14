package models

import "net"

type Connection struct {
	Conn net.Conn
}

type Topic string

type Subscription struct {
	Topic string
	Conn  Connection
}
