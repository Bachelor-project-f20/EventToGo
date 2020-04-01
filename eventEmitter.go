package eventToGo

import models "github.com/Bachelor-project-f20/eventToGo/shared/proto/gen"

type EventEmitter interface {
	Emit(e models.Event) error
}
