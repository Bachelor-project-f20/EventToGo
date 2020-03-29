package nats

import (
	"log"

	etg "github.com/Bachelor-project-f20/eventToGo"
	"github.com/nats-io/go-nats"
)

type natsEventEmitter struct {
	connection *nats.EncodedConn
	exchange   string
	queue      string
}

func NewNatsEventEmitter(connection *nats.EncodedConn, exchange, queue string) (etg.EventEmitter, error) {
	emitter := natsEventEmitter{
		connection,
		exchange,
		queue,
	}

	return &emitter, nil
}

func (n *natsEventEmitter) Emit(e etg.BaseEvent) error {
	err := n.connection.Publish(e.GetEventName(), e)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
