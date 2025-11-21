package promotion

import 	"github.com/alan-b-lima/almodon/pkg/uuid"


type Service interface {
	List(ListParams) (Entities, error)

	Get(uuid.UUID) (Entity, error)
	GetByUser(uuid.UUID) (Entity, error)

	Create(Create) (uuid.UUID, error)

	Update(uuid.UUID, Update) error

	Delete(uuid.UUID) error
}
