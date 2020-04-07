package nats_test

import (
	"fmt"
	"testing"
	"time"

	lnats "github.com/Bachelor-project-f20/eventToGo/nats"
	models "github.com/Bachelor-project-f20/shared/models"
)

var connectionString = "localhost:4222"
var exchange = "test"
var queueType = "queue"
var natsHandler *lnats.NatsHandler

func TestEmit(t *testing.T) {
	eventEmitter, err := natsHandler.SetupEmitter(exchange, queueType, connectionString)

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
	eventListener, err := natsHandler.SetupListener(exchange, queueType, connectionString)

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
	eventEmitter, _ := natsHandler.SetupEmitter(exchange, queueType, connectionString)
	eventEmitter.Emit(event)

	recEvent := <-eventChan

	fmt.Printf("Event received in test: %v", recEvent)
}

func TestCreateEmitterAndListener(t *testing.T) {
	_, _, err := natsHandler.SetupEmitterAndListener(exchange, queueType, connectionString)

	if err != nil {
		fmt.Println("Error while creating emitter and listener simultaneously")
		t.Error(err)
	}
}
