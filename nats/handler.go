package nats

import (
	"fmt"
	"log"
	"time"

	etg "github.com/Bachelor-project-f20/eventToGo"
	"github.com/nats-io/nats.go"
)

type NatsHandler struct{}

func (n *NatsHandler) SetupEmitter(exchange, queueType, connectionString string) (eventEmitter etg.EventEmitter, err error) {
	encodedConn, err := setupNatsConn(connectionString)
	return setupEmitter(exchange, queueType, encodedConn)
}

func setupEmitter(exchange, queueType string, encodedConn *nats.EncodedConn) (eventEmitter etg.EventEmitter, err error) {
	eventEmitter, err = newNatsEventEmitter(encodedConn, exchange, queueType)
	if err != nil {
		log.Fatalf("Error creating event emitter: %v \n", err)
		return nil, err
	}
	log.Println("EventToGo: nats emitter done")
	return eventEmitter, nil
}

func (n *NatsHandler) SetupListener(exchange, queueType, connectionString string) (eventListener etg.EventListener, err error) {
	encodedConn, err := setupNatsConn(connectionString)
	return setupListener(exchange, queueType, encodedConn)
}

func setupListener(exchange, queueType string, encodedConn *nats.EncodedConn) (eventListener etg.EventListener, err error) {
	eventListener, err = newNatsEventListener(encodedConn, exchange, queueType)
	if err != nil {
		log.Fatalf("Error creating event listener: %v \n", err)
		return nil, err
	}
	log.Println("EventToGo: nats listener done")
	return eventListener, nil
}

func (n *NatsHandler) SetupEmitterAndListener(exchange, queueType, connectionString string) (eventEmitter etg.EventEmitter, eventListener etg.EventListener, err error) {
	encodedConn, err := setupNatsConn(connectionString)
	eventEmitter, err = setupEmitter(exchange, queueType, encodedConn)
	if err != nil {
		return nil, nil, err
	}
	eventListener, err = setupListener(exchange, queueType, encodedConn)
	if err != nil {
		return nil, nil, err
	}
	return eventEmitter, eventListener, nil
}

func setupNatsConn(natsConnectionString string) (*nats.EncodedConn, error) {

	natsConn := connect(natsConnectionString)

	encodedConn, err := nats.NewEncodedConn(natsConn, "json")

	if err != nil {
		log.Fatalln("Creation of encoded connection failed")
		return nil, err
	}

	return encodedConn, nil

}

func connect(natsConnectionString string) *nats.Conn {
	i := 5
	for i > 0 {
		natsConn, err := nats.Connect(natsConnectionString)
		if err != nil {
			fmt.Println("Can't connect to nats, sleeping for 2 sec, err: ", err)
			time.Sleep(2 * time.Second)
			i--
			continue
		} else {
			fmt.Println("Connected to nats")
			return natsConn
		}
	}
	panic("Not connected to nats")
}
