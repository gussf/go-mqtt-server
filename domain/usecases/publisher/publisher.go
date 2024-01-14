package publisher

import (
	"errors"
	"fmt"
	"slices"

	"github.com/gussf/go-mqtt-server/domain/models"
)

type Usecase interface {
	AddConnectionToTopicPool(models.Subscription) error
	PublishToSubscribers(models.Subscription, []byte) error
	IsSubscribed(models.Subscription) bool
}

type usecase struct {
	pool models.SubscriptionPool
}

func NewUsecase() Usecase {
	return usecase{
		pool: models.NewSubscriptionPool(),
	}
}

func (uc usecase) AddConnectionToTopicPool(sub models.Subscription) error {
	if sub.Topic == "" {
		return errors.New("invalid topic name")
	}

	if sub.Conn.Conn == nil {
		return errors.New("invalid connection")
	}

	uc.pool.Add(sub)

	fmt.Printf("Subscribed connection %v to topic %s\n", sub.Conn, sub.Topic)
	return nil
}

func (uc usecase) PublishToSubscribers(sub models.Subscription, payload []byte) error {
	conns, ok := uc.pool.Get(sub.Topic)
	if !ok {
		return errors.New("failed to find topic in pool")
	}

	for _, c := range conns {
		c := c
		go func() {
			n, err := c.Conn.Write(payload)
			fmt.Printf("writing %x to %v\n", payload, c.Conn)
			fmt.Printf("wrote: %d\nerr: %v\n", n, err)
		}()
	}
	return nil
}

func (uc usecase) IsSubscribed(sub models.Subscription) bool {
	conns, ok := uc.pool.Get(sub.Topic)
	if !ok {
		return false
	}

	if slices.Contains(conns, sub.Conn) {
		return true
	}

	return false
}
