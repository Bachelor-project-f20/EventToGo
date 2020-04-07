package sns_test

// To run the tests, a local Docker container, based on the following Docker image
// must be running, on port 9911 (unless the test code is changed):
// https://hub.docker.com/r/s12v/sns/?fbclid=IwAR23X1mEVHH5Q64awf-ZtyzC_r712-yjfmqEQGRvDCT8LYfMkdyP4goTxdE

//Alternatively, one can attach to an actual SNS instance, by using the SharedConfigState session initialization

import (
	"fmt"
	"testing"
	"time"

	lsns "github.com/Bachelor-project-f20/eventToGo/sns"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var svc *sns.SNS
var topicArnMap map[string]string
var events []string

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

	//creating subscriptions now, since AWS doesn't allow publishing to subjects without subscribers
	for count, _ := range events {
		event := events[count]
		arn := topicArnMap[event]
		_, err := svc.Subscribe(&sns.SubscribeInput{
			//must provide the acutal IP of your machine, does not work with localhost - not important for this test however
			Endpoint:              aws.String("http://<LOCAL-IP-HERE>/snstest"), //Get actual IP address somehow - OS package?
			Protocol:              aws.String("http"),
			ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
			TopicArn:              aws.String(arn),
		})
		if err != nil {
			fmt.Printf("Error subscribing to SNS topic: %v \n", err)
			t.Error(err)
		}
	}

	emitErr := emitter.Emit(event)

	if emitErr != nil {
		fmt.Println("Error while emitting event")
		t.Error(err)
	}
	fmt.Println("Event emitted")
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
		t.Error(listenErr)
	}

	emitter, _ := lsns.NewSNSEventEmitter(svc, topicArnMap)
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
