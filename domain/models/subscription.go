package models

type SubscriptionPool struct {
	subscriptions map[Topic][]Connection
}
