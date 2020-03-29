package eventToGo

type EventEmitter interface {
	Emit(e BaseEvent) error
}
