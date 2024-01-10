package publisher_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gussf/go-mqtt-server/domain/models"
	"github.com/gussf/go-mqtt-server/domain/usecases/publisher"
)

func TestAddConnectionToTopicPool(t *testing.T) {
	t.Run("", func(t *testing.T) {
		p := publisher.NewUsecase()

		sub := models.Subscription{
			Conn:  models.Connection{Conn: &net.IPConn{}},
			Topic: models.Topic{Name: "test"},
		}
		err := p.AddConnectionToTopicPool(sub.Conn, sub.Topic)
		assert.Nil(t, err)
		err = p.AddConnectionToTopicPool(sub.Conn, sub.Topic)
		assert.NotNil(t, err)
		// wip
	})
}
