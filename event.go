package eventToGo

type BaseEvent interface {
	GetID() string
	GetEventName() string
	GetPublisher() string
	GetTimestamp() int64
	GetPayload() []byte
}

type Event struct {
	ID        string
	EventName string
	Publisher string
	TimeStamp int64
	Payload   []byte
}

func (e Event) GetID() string {
	return e.ID
}

func (e Event) GetEventName() string {
	return e.EventName
}

func (e Event) GetPublisher() string {
	return e.Publisher
}

func (e Event) GetTimestamp() int64 {
	return e.TimeStamp
}

func (e Event) GetPayload() []byte {
	return e.Payload
}
