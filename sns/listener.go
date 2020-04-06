package sns

import (
	"encoding/base64"
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

	err := s.createSubscriptions(events...)

	if err != nil {
		log.Fatalf("Error creating subscriptions: %v \n", err)
		return nil, nil, err
	}

	createSubscriptionHandlers(eventChan)
	return eventChan, errChan, nil
}

func (s *snsEventListener) createSubscriptions(events ...string) error {
	for count, _ := range events {
		event := events[count]
		arn := s.topicArnMap[event]
		result, err := s.client.Subscribe(&sns.SubscribeInput{
			Endpoint:              aws.String("http://localhost:8081/"), //Get actual IP address somehow - OS package?
			Protocol:              aws.String("http"),
			ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
			TopicArn:              aws.String(arn),
		})
		if err != nil {
			log.Fatalf("Error subscribing to SNS topic: %v \n", err)
			return err
		}
		fmt.Printf("Created subscription: %s \n", *result.SubscriptionArn)
	}
	return nil
}

//Everything below was shamelessly repurposed from: https://github.com/viveksyngh/aws-sns-subscriber/blob/master/subscriber/subscriber.go

const subConfrmType = "SubscriptionConfirmation"
const notificationType = "Notification"

func createSubscriptionHandlers(eventChan chan models.Event) {
	handler := chanHandler(eventChan)
	http.HandleFunc("/", handler)
	go http.ListenAndServe(":8081", nil)
}

//handler processes messages sent by SNS
func chanHandler(eventChan chan<- models.Event) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var f interface{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Unable to Parse Body: %v \n", err)
			return
		}
		fmt.Printf(string(body))
		err = json.Unmarshal(body, &f)
		if err != nil {
			fmt.Printf("Unable to Unmarshal request: %v \n", err)
			return
		}

		data := f.(map[string]interface{})
		fmt.Println(data["Type"].(string))

		if data["Type"].(string) == subConfrmType {
			subcribeURL := data["SubscribeURL"].(string)
			go confirmSubscription(subcribeURL)
		} else if data["Type"].(string) == notificationType {

			message := data["Message"]
			byteMessage, _ := json.Marshal(message) //So we're json marshalling this, to unmarshal it again, which is the dumbest shit ever, but it works

			var event models.Event
			err := json.Unmarshal(byteMessage, &event)

			if err != nil {
				fmt.Printf("Unable to Unmarshal request: %v \n", err)
				return
			}

			payload := event.Payload
			var decodedPayload []byte
			base64.StdEncoding.Decode(decodedPayload, payload)
			event.Payload = decodedPayload

			eventChan <- event
		}
	}
}

func confirmSubscription(subscriptionUrl string) error {
	response, err := http.Get(subscriptionUrl)
	if err != nil {
		fmt.Printf("Unable to confirm subscriptions")
		return err
	} else {
		fmt.Printf("Subscription Confirmed sucessfully. %d", response.StatusCode)
	}
	return nil
}
