package publisher_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gussf/go-mqtt-server/domain/models"
	"github.com/gussf/go-mqtt-server/domain/usecases/publisher"
)

func TestAddConnectionToTopicPool(t *testing.T) {

	t.Run("should add connection succesfully", func(t *testing.T) {
		s, c := net.Pipe()
		defer c.Close()
		defer s.Close()

		p := publisher.NewUsecase()

		sub := models.Subscription{
			Conn:  models.Connection{Conn: c},
			Topic: "test",
		}

		err := p.AddConnectionToTopicPool(sub)
		assert.Nil(t, err)

		subbed := p.IsSubscribed(sub)
		assert.True(t, subbed)
	})
}
