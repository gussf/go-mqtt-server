package models

import (
	"slices"
)

type SubscriptionPool interface {
	Add(sub Subscription)
	Get() (Subscription, bool)
}

type subscriptionPool struct {
	topicPool map[Topic][]Connection
}

func NewSubscriptionPool() SubscriptionPool {
	return subscriptionPool{
		topicPool: map[Topic][]Connection{},
	}
}

func (s subscriptionPool) Add(sub Subscription) {
	pool := s.topicPool[sub.Topic]

	if slices.Contains(pool, sub.Conn) {
		return
	}

	pool = append(pool, sub.Conn)
}

func (s subscriptionPool) Get() (Subscription, bool) {
	return Subscription{}, false
}
