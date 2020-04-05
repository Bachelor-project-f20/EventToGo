package sns_test

import (
	"fmt"
	"testing"
	"time"

	lsns "github.com/Bachelor-project-f20/eventToGo/sns"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var svc *sns.SNS
var topicArnMap map[string]string
var events []string

func TestSetupSNS(t *testing.T) {

	//Session instantation currently failing, not sure what to put in to use the s12v/sns container
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Endpoint: aws.String("http://localhost:9911"), Region: aws.String("us-east-1"), DisableSSL: aws.Bool(true)},
	}))
	svc = sns.New(sess)

	events = []string{
		"test",
	}
	var err error
	topicArnMap, err = setupTopicArns(events)

	if err != nil {
		fmt.Printf("Error creating topics with SNS instance, error: %v \n", err)
		t.Error(err)
	}

}

func TestEmit(t *testing.T) {

	emitter, err := lsns.NewSNSEventEmitter(svc, topicArnMap)

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
	fmt.Println("Event emited")

}

func TestListen(t *testing.T) {

	listener, err := lsns.NewSNSEventListener(svc, topicArnMap)

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
		t.Error(err)
	}

	emitter, err := lsns.NewSNSEventEmitter(svc, topicArnMap)
	emitter.Emit(event)

	recEvent := <-eventChan

	fmt.Printf("Event received in test: %v", recEvent)

}

func setupTopicArns(events []string) (map[string]string, error) {
	arnMap := make(map[string]string)
	for count, _ := range events {
		eventName := events[count]
		result, err := svc.CreateTopic(&sns.CreateTopicInput{
			Name: aws.String(eventName),
		})

		if err != nil {
			fmt.Printf("Error creating topics with SNS instance, error: %v \n", err)
			return nil, err
		}

		arnMap[eventName] = *result.TopicArn
	}
	return arnMap, nil
}
