package nats

import (
	etg "github.com/Bachelor-project-f20/eventToGo"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/nats-io/nats.go"
)

type natsEventListener struct {
	connection *nats.EncodedConn
	exchange   string
	queue      string
}

func newNatsEventListener(connection *nats.EncodedConn, exchange, queue string) (etg.EventListener, error) {
	listener := natsEventListener{
		connection,
		exchange,
		queue,
	}
	return &listener, nil
}

func (n *natsEventListener) Listen(events ...string) (<-chan models.Event, <-chan error, error) {
	eventChan := make(chan models.Event)
	errChan := make(chan error)

	for count, _ := range events {
		_, err := n.connection.QueueSubscribe(events[count], n.queue, func(e *models.Event) {
			eventChan <- *e
		})
		if err != nil {
			errChan <- err
		}
	}
	return eventChan, errChan, nil
}
