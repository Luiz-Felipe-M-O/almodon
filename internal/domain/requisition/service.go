package requisition

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	List(req ListParams) (Entities, error)
	Get(uuid uuid.UUID) (Entity, error)

	Create(req Create) (uuid.UUID, error)

	Patch(uuid uuid.UUID, req Patch) error
	Delete(uuid uuid.UUID) error

	AddEntry(requisitionUUID uuid.UUID, req AddEntry) (uuid.UUID, error)
	RemoveEntry(requisitionUUID, entryUUID uuid.UUID) error

	Answer(requisitionUUID uuid.UUID, req AnswerRequisition) error
	Cancel(requisitionUUID uuid.UUID) error
	MarkFulfilled(requisitionUUID uuid.UUID) error

	Allow(auth.Actor) Service
}
