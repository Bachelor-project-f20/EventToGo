package nats_test

import (
	"fmt"
	"testing"
	"time"

	lnats "github.com/Bachelor-project-f20/eventToGo/nats"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/nats-io/go-nats"
)

var natsConn *nats.Conn
var encodedConn *nats.EncodedConn
var exchange string
var queueType string

func TestSetup(t *testing.T) {
	natsConn, err := nats.Connect("localhost:4222")

	if err != nil {
		fmt.Println("Connection to Nats failed")
		t.Error(err)
	}

	encodedConn, err = nats.NewEncodedConn(natsConn, "json")

	if err != nil {
		fmt.Println("Creation of encoded connection failed")
		t.Error(err)
	}

	exchange = "test"
	queueType = "queue"

}

func TestEmit(t *testing.T) {
	eventEmitter, err := lnats.NewNatsEventEmitter(encodedConn, exchange, queueType)

	if err != nil {
		fmt.Println("Creation of Emitter  failed")
		t.Error(err)
	}

	event := models.Event{}
	event.ID = "test"
	event.EventName = "test"
	event.Publisher = "test"
	event.Timestamp = time.Now().UnixNano()
	event.Payload = []byte{'t'}

	emitErr := eventEmitter.Emit(event)

	if emitErr != nil {
		fmt.Println("Error while emitting event")
		t.Error(err)
	}
	fmt.Println("Event emited")
}

func TestListen(t *testing.T) {
	eventListener, err := lnats.NewNatsEventListener(encodedConn, exchange, queueType)

	if err != nil {
		fmt.Println("Creation of Listener  failed")
		t.Error(err)
	}

	event := models.Event{}
	event.ID = "test"
	event.EventName = "test"
	event.Publisher = "test"
	event.Timestamp = time.Now().UnixNano()
	event.Payload = []byte{'t'}

	eventChan, _, listenErr := eventListener.Listen(event.ID)
	if listenErr != nil {
		fmt.Println("Listen function  failed")
		t.Error(err)
	}

	//Necessary - when the Nats connection is not set to durable, messages in unsubscribed message queues are lost
	eventEmitter, _ := lnats.NewNatsEventEmitter(encodedConn, exchange, queueType)
	eventEmitter.Emit(event)

	recEvent := <-eventChan

	fmt.Printf("Event received in test: %v", recEvent)
}
