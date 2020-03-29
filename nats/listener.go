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

type Event struct {
	ID        string
	EventName string
	Publisher string
	Timestamp int64
	Payload   []byte
}

func (e *Event) GetID() string {
	return e.ID
}

func (e *Event) GetPublisher() string {
	return e.Publisher
}

func (e *Event) GetTimestamp() int64 {
	return e.Timestamp
}

func (e *Event) GetEventName() string {
	return e.EventName
}

func (e *Event) GetPayload() []byte {
	return e.Payload
}

func NewNatsEventListener(connection *nats.EncodedConn, exchange, queue string) (etg.EventListener, error) {
	listener := natsEventListener{
		connection,
		exchange,
		queue,
	}
	return &listener, nil
}

func (n *natsEventListener) Listen(events ...string) (<-chan etg.BaseEvent, <-chan error, error) {
	eventChan := make(chan etg.BaseEvent)
	errChan := make(chan error)

	for count, _ := range events {
		_, err := n.connection.QueueSubscribe(events[count], n.queue, func(e *Event) {
			eventChan <- e
		})
		if err != nil {
			errChan <- err
		}
	}
	return eventChan, errChan, nil
}
