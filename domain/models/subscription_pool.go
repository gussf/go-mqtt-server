package models

import (
	"slices"
)

type SubscriptionPool interface {
	Add(Subscription)
	Get(string) ([]Connection, bool)
}

type subscriptionPool struct {
	topicPool map[string][]Connection
}

func NewSubscriptionPool() SubscriptionPool {
	return subscriptionPool{
		topicPool: map[string][]Connection{},
	}
}

func (s subscriptionPool) Add(sub Subscription) {
	pool := s.topicPool[sub.Topic]

	if slices.Contains(pool, sub.Conn) {
		return
	}

	s.topicPool[sub.Topic] = append(pool, sub.Conn)
}

func (s subscriptionPool) Get(topic string) ([]Connection, bool) {
	connSlice, ok := s.topicPool[topic]
	if !ok {
		return nil, false
	}

	return connSlice, true
}
