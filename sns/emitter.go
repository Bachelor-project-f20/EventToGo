package sns

import (
	"encoding/json"
	"log"

	etg "github.com/Bachelor-project-f20/eventToGo"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type snsEventEmitter struct {
	client      *sns.SNS
	topicArnMap map[string]string
}

func newSNSEventEmitter(client *sns.SNS, topicArnMap map[string]string) (etg.EventEmitter, error) {
	return &snsEventEmitter{client, topicArnMap}, nil
}

func (s *snsEventEmitter) Emit(event models.Event) error {

	marshalEvent, err := json.Marshal(event)

	if err != nil {
		log.Fatalf("Error parsing event to json: %v \n", err)
		return err
	}

	stringMarshalEvent := string(marshalEvent)

	input := &sns.PublishInput{
		Message:  aws.String(stringMarshalEvent),
		TopicArn: aws.String(s.topicArnMap[event.EventName]),
	}

	_, err = s.client.Publish(input)
	if err != nil {
		log.Fatalf("Publish error: %v \n", err)
		return err
	}

	return nil
}
