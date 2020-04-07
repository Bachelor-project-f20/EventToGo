package sns_test

// To run the tests, a local Docker container, based on the following Docker image
// must be running, on port 9911 (unless the test code is changed):
// https://hub.docker.com/r/s12v/sns/?fbclid=IwAR23X1mEVHH5Q64awf-ZtyzC_r712-yjfmqEQGRvDCT8LYfMkdyP4goTxdE

//Alternatively, one can attach to an actual SNS instance, by using the SharedConfigState session initialization

import (
	"fmt"
	"testing"
	"time"

	handler "github.com/Bachelor-project-f20/eventToGo/sns"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var svc *sns.SNS
var events []string
var snsHandler *handler.SNSHandler

func TestSetupSNS(t *testing.T) {

	//AnonymousCredentials for the mock SNS instance
	//SSL disabled, because it's easier when testing
	//localhost:991 is where the fake SNS container should be running
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Credentials: credentials.AnonymousCredentials, Endpoint: aws.String("http://localhost:9911"), Region: aws.String("us-east-1"), DisableSSL: aws.Bool(true)},
	}))

	//Using actual SNS instance
	//Requires config and credential file in ~/.aws/
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))
	svc = sns.New(sess)

	events = []string{
		"test",
	}

}

func TestListen(t *testing.T) {

	listener, err := snsHandler.SetupListener(svc, events...)

	if err != nil {
		fmt.Printf("Error creating listener, error: %v \n", err)
		t.Error(err)
	}

	event := models.Event{}
	event.ID = "test"
	event.EventName = "test"
	event.Publisher = "test"
	event.Timestamp = time.Now().UnixNano()
	event.Payload = []byte{'t'}

	eventChan, _, listenErr := listener.Listen(event.EventName)

	if listenErr != nil {
		fmt.Println("Listen function  failed")
		t.Error(listenErr)
	}

	emitter, _ := snsHandler.SetupEmitter(svc, events...)
	emitter.Emit(event)

	recEvent := <-eventChan

	fmt.Printf("Event received in test: %v", recEvent)

}

func TestEmit(t *testing.T) {

	emitter, err := snsHandler.SetupEmitter(svc, events...)

	if err != nil {
		fmt.Printf("Error creating emitter, error: %v \n", err)
		t.Error(err)
	}

	event := models.Event{}
	event.ID = "test"
	event.EventName = "test"
	event.Publisher = "test"
	event.Timestamp = time.Now().UnixNano()
	event.Payload = []byte{'t'}

	emitErr := emitter.Emit(event)

	if emitErr != nil {
		fmt.Println("Error while emitting event")
		t.Error(err)
	}
	fmt.Println("Event emitted")
}

func TestCreateEmitterAndListener(t *testing.T) {
	_, _, err := snsHandler.SetupEmitterAndListener(svc, events...)

	if err != nil {
		fmt.Println("Error while creating emitter and listener simultaneously")
		t.Error(err)
	}

}
