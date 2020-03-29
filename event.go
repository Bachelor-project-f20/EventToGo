package eventToGo

type Event struct {
	ID        string
	EventName string
	Publisher string
	Timestamp int64
	Payload   []byte
}
