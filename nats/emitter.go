package nats

import (
	"log"

	etg "github.com/Bachelor-project-f20/eventToGo"
	models "github.com/Bachelor-project-f20/shared/models"
	nats "github.com/nats-io/nats.go"
)

type natsEventEmitter struct {
	connection *nats.EncodedConn
	exchange   string
	queue      string
}

func newNatsEventEmitter(connection *nats.EncodedConn, exchange, queue string) (etg.EventEmitter, error) {
	emitter := natsEventEmitter{
		connection,
		exchange,
		queue,
	}

	return &emitter, nil
}

func (n *natsEventEmitter) Emit(e models.Event) error {
	err := n.connection.Publish(e.EventName, e)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
