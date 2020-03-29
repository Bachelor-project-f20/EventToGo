package eventToGo

type EventEmitter interface {
	Emit(e Event) error
}
