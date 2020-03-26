package eventToGo

type Event interface {
	GetID() string
	GetPublisher() string
	GetTimestamp() int64
}
