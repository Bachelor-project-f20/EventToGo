package eventToGo

type EventListener interface {
	Listen(events ...string) (<-chan BaseEvent, <-chan error, error)
}
