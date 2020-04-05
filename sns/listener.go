package sns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	etg "github.com/Bachelor-project-f20/eventToGo"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type snsEventListener struct {
	client      *sns.SNS
	topicArnMap map[string]string
}

func NewSNSEventListener(client *sns.SNS, topicArnMap map[string]string) (etg.EventListener, error) {
	return &snsEventListener{client, topicArnMap}, nil
}

func (s *snsEventListener) Listen(events ...string) (<-chan models.Event, <-chan error, error) {
	eventChan := make(chan models.Event)
	errChan := make(chan error)

	for count, _ := range events {
		event := events[count]
		arn := s.topicArnMap[event]
		result, err := s.client.Subscribe(&sns.SubscribeInput{
			Endpoint:              aws.String("http://localhost:8081"), //Get actual IP address somehow - OS package?
			Protocol:              aws.String("http"),
			ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
			TopicArn:              aws.String(arn),
		})
		if err != nil {
			log.Fatalf("Error subscribing to SNS topic: %v \n", err)
			return nil, nil, err
		}

		fmt.Println(*result.SubscriptionArn)
	}

	createSubscriptionHandlers(eventChan)
	return eventChan, errChan, nil
}

//Everything below was shamelessly repurposed from: https://github.com/viveksyngh/aws-sns-subscriber/blob/master/subscriber/subscriber.go

func createSubscriptionHandlers(eventChan chan models.Event) {
	handler := chanHandler(eventChan)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}

//handler processes messages sent by SNS
func chanHandler(eventChan chan models.Event) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Unable to Parse Body")
		}
		fmt.Printf(string(body))
		var event models.Event
		err = json.Unmarshal(body, &event)
		if err != nil {
			fmt.Printf("Unable to Unmarshal request")
		}
		eventChan <- event
	}
}
