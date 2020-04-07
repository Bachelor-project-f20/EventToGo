package sns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	etg "github.com/Bachelor-project-f20/eventToGo"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type snsEventListener struct {
	client      *sns.SNS
	topicArnMap map[string]string
}

func newSNSEventListener(client *sns.SNS, topicArnMap map[string]string) (etg.EventListener, error) {
	return &snsEventListener{client, topicArnMap}, nil
}

func (s *snsEventListener) Listen(events ...string) (<-chan models.Event, <-chan error, error) {
	eventChan := make(chan models.Event)
	errChan := make(chan error)

	err := s.createSubscriptions(events...)
	createSubscriptionHandlers(eventChan)

	if err != nil {
		log.Fatalf("Error creating subscriptions: %v \n", err)
		return nil, nil, err
	}

	return eventChan, errChan, nil
}

func (s *snsEventListener) createSubscriptions(events ...string) error {
	for count, _ := range events {
		event := events[count]
		arn := s.topicArnMap[event]
		endpoint := "http://" + getLocalIP() + ":8081/sns"
		_, err := s.client.Subscribe(&sns.SubscribeInput{
			//must provide the acutal IP of your machine, does not work with localhost
			Endpoint:              aws.String(endpoint), //Get actual IP address somehow - OS package?
			Protocol:              aws.String("http"),
			ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
			TopicArn:              aws.String(arn),
		})
		if err != nil {
			log.Fatalf("Error subscribing to SNS topic: %v \n", err)
			return err
		}
	}
	return nil
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		panic(err)
	}

	defer conn.Close()
	localAddr := conn.LocalAddr()
	ipAddr := strings.Split(localAddr.String(), ":")
	return ipAddr[0]
}

//Everything below was shamelessly repurposed from:
// https://github.com/viveksyngh/aws-sns-subscriber/blob/master/subscriber/subscriber.go

const subConfrmType = "SubscriptionConfirmation"
const notificationType = "Notification"

func createSubscriptionHandlers(eventChan chan models.Event) {
	handler := chanHandler(eventChan)
	http.HandleFunc("/sns", handler)
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
		err = json.Unmarshal(body, &f)
		if err != nil {
			fmt.Printf("Unable to Unmarshal request: %v \n", err)
			return
		}

		data := f.(map[string]interface{})

		if data["Type"].(string) == subConfrmType {
			subcribeURL := data["SubscribeURL"].(string)
			go confirmSubscription(subcribeURL)
		} else if data["Type"].(string) == notificationType {

			message := data["Message"].(string)
			byteMessage := []byte(message)

			if err != nil {
				log.Fatalf("Unable to marshal incoming message into byte array: %v \n", err)
				return
			}

			var event models.Event
			err = json.Unmarshal(byteMessage, &event)

			if err != nil {
				fmt.Printf("Unable to Unmarshal request: %v \n", err)
				return
			}

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
