package sns

import (
	"fmt"
	"log"

	etg "github.com/Bachelor-project-f20/eventToGo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSHandler struct{}

func (s *SNSHandler) SetupEmitter(client *sns.SNS, events ...string) (eventEmitter etg.EventEmitter, err error) {
	topicArnMap, err := setupTopicArns(client, events)

	if err != nil {
		log.Fatalf("Error creating topic ARNs: %v \n", err)
		return nil, err
	}

	return setupEmitter(client, topicArnMap)
}
func setupEmitter(client *sns.SNS, topicArnMap map[string]string) (eventEmitter etg.EventEmitter, err error) {
	emitter, err := newSNSEventEmitter(client, topicArnMap)
	if err != nil {
		log.Fatalf("Error creating event emitter: %v \n", err)
		return nil, err
	}
	log.Println("EventToGo: SNS emitter done")
	return emitter, nil
}

func (s *SNSHandler) SetupListener(client *sns.SNS, events ...string) (eventListener etg.EventListener, err error) {
	topicArnMap, err := setupTopicArns(client, events)

	if err != nil {
		log.Fatalf("Error creating topic ARNs: %v \n", err)
		return nil, err
	}

	return setupListener(client, topicArnMap)
}

func setupListener(client *sns.SNS, topicArnMap map[string]string) (eventListener etg.EventListener, err error) {
	listener, err := newSNSEventListener(client, topicArnMap)
	if err != nil {
		log.Fatalf("Error creating event listener: %v \n", err)
		return nil, err
	}
	log.Println("EventToGo: SNS listener done")
	return listener, nil
}

func (s *SNSHandler) SetupEmitterAndListener(client *sns.SNS, events ...string) (eventEmitter etg.EventEmitter, eventListener etg.EventListener, err error) {
	topicArnMap, err := setupTopicArns(client, events)

	if err != nil {
		log.Fatalf("Error creating topic ARNs: %v \n", err)
		return nil, nil, err
	}

	listener, err := newSNSEventListener(client, topicArnMap)
	if err != nil {
		log.Fatalf("Error creating event listener: %v \n", err)
		return nil, nil, err
	}
	log.Println("EventToGo: SNS listener done")

	emitter, err := newSNSEventEmitter(client, topicArnMap)
	if err != nil {
		log.Fatalf("Error creating event emitter: %v \n", err)
		return nil, nil, err
	}
	log.Println("EventToGo: SNS emitter done")
	return emitter, listener, nil
}

func setupTopicArns(client *sns.SNS, events []string) (map[string]string, error) {
	arnMap := make(map[string]string)
	for count, _ := range events {
		eventName := events[count]
		result, err := client.CreateTopic(&sns.CreateTopicInput{
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
