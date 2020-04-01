package eventToGo

import models "github.com/Bachelor-project-f20/shared/models"

type EventListener interface {
	Listen(events ...string) (<-chan models.Event, <-chan error, error)
}
