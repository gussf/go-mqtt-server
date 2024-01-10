package models

import "net"

type Connection struct {
	Conn net.Conn
}

type Topic struct {
	Name string
}

type Subscription struct {
	Topic Topic
	Conn  Connection
}

