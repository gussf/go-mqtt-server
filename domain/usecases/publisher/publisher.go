package publisher

import (
	"errors"

	"github.com/gussf/go-mqtt-server/domain/models"
)

type Usecase interface {
	AddConnectionToTopicPool(models.Connection, models.Topic) error
	PublishToSubscribers(models.Topic) error
}

type usecase struct {
	pool models.SubscriptionPool
}

func NewUsecase() Usecase {
	return usecase{
		pool: models.NewSubscriptionPool(),
	}
}

func (uc usecase) AddConnectionToTopicPool(conn models.Connection, topic models.Topic) error {
	if topic.Name == "" {
		return errors.New("invalid topic name")
	}

	sub := models.Subscription{
		Topic: topic,
		Conn:  conn,
	}

	uc.pool.Add(sub)

	return nil
}

func (uc usecase) PublishToSubscribers(topic models.Topic) error {
	return nil
}
