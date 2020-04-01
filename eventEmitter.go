package eventToGo

import models "github.com/Bachelor-project-f20/shared/models"

type EventEmitter interface {
	Emit(e models.Event) error
}
