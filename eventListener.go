package eventToGo

import models "github.com/Bachelor-project-f20/eventToGo/shared/proto/gen"

type EventListener interface {
	Listen(events ...string) (<-chan models.Event, <-chan error, error)
}
