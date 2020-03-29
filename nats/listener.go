package nats

import (
	etg "github.com/Bachelor-project-f20/eventToGo"
	"github.com/nats-io/go-nats"
)

type natsEventListener struct {
	connection *nats.EncodedConn
	exchange   string
	queue      string
}

func NewNatsEventListener(connection *nats.EncodedConn, exchange, queue string) (etg.EventListener, error) {
	listener := natsEventListener{
		connection,
		exchange,
		queue,
	}
	return &listener, nil
}

func (n *natsEventListener) Listen(events ...string) (<-chan etg.Event, <-chan error, error) {
	eventChan := make(chan etg.Event)
	errChan := make(chan error)

	for count, _ := range events {
		_, err := n.connection.QueueSubscribe(events[count], n.queue, func(e *etg.Event) {
			eventChan <- *e
		})
		if err != nil {
			errChan <- err
		}
	}
	return eventChan, errChan, nil
}
