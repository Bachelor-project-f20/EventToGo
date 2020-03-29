package eventToGo

type Event struct {
	ID        string
	EventName string
	Publisher string
	TimeStamp int64
	Payload   []byte
}
