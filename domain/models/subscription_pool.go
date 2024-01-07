package models

type SubscriptionPool interface {
	Add() error
}

type subscriptionPool struct {
	topicPool map[Topic][]Connection
}

func NewSubscriptionPool() SubscriptionPool {
	return &subscriptionPool{
		topicPool: map[Topic][]Connection{},
	}
}

func (s subscriptionPool) Add() error {
	return nil
}
