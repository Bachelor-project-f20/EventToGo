package eventToGo

type Event interface {
	GetID() string
	GetEventName() string
	GetPublisher() string
	GetTimestamp() int64
	GetPayload() []byte
}
