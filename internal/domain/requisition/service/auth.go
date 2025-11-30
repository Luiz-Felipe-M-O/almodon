package requisitionserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/requisition"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	requisition.Service
	actor auth.Actor
}

func New(service requisition.Service) requisition.Service {
	return &Gate{Service: service}
}

var (
	permLogged = auth.Permit(auth.User, auth.Admin, auth.Promoted, auth.Chief)

	permAdmin = auth.Permit(auth.Admin, auth.Promoted, auth.Chief)

	permPromoted = auth.Permit(auth.Promoted, auth.Chief)
)

func (s *Gate) Allow(act auth.Actor) requisition.Service {
	return &Gate{
		Service: s.Service,
		actor:   act,
	}
}

func (s *Gate) List(req requisition.ListParams) (requisition.Entities, error) {
	if err := service.Authorize(permLogged, s.actor); err != nil {
		return requisition.Entities{}, err
	}

	if s.actor.Role() == auth.User {
		req.Author = opt.Some(s.actor.User().String())
	}

	return s.Service.List(req)
}

func (s *Gate) Get(uuid uuid.UUID) (requisition.Entity, error) {
	if err := service.Authorize(permLogged, s.actor); err != nil {
		return requisition.Entity{}, err
	}

	res, err := s.Service.Get(uuid)
	if err != nil {
		return requisition.Entity{}, err
	}

	if s.actor.Role() == auth.User && res.Author != s.actor.User() {
		return requisition.Entity{}, service.ErrUnauthorized
	}

	return res, nil
}

func (s *Gate) Create(req requisition.Create) (uuid.UUID, error) {
	if err := service.Authorize(permLogged, s.actor); err != nil {
		return uuid.UUID{}, err
	}

	r, err := requisition.New(
		s.actor.User(),
		req.Notes,
		req.Destination,
		req.Entries,
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	return r.UUID(), s.Service.Create(requisition.Create{
		Notes:       r.Notes(),
		Destination: r.Destination(),
		Entries:     req.Entries,
	})
}

func (s *Gate) Patch(uuid uuid.UUID, req requisition.Patch) error {
	if err := service.Authorize(permLogged, s.actor); err != nil {
		return err
	}

	res, err := s.Service.Get(uuid)
	if err != nil {
		return err
	}

	if s.actor.Role() == auth.User && res.Author != s.actor.User() {
		return service.ErrUnauthorized
	}

	return s.Service.Patch(uuid, req)
}

func (s *Gate) Delete(uuid uuid.UUID) error {
	if err := service.Authorize(permLogged, s.actor); err != nil {
		return err
	}

	res, err := s.Service.Get(uuid)
	if err != nil {
		return err
	}

	if s.actor.Role() == auth.User && res.Author != s.actor.User() {
		return service.ErrUnauthorized
	}

	return s.Service.Delete(uuid)
}

func (s *Gate) AddEntry(requisitionUUID uuid.UUID, req requisition.AddEntry) (uuid.UUID, error) {
	if err := service.Authorize(permLogged, s.actor); err != nil {
		return uuid.UUID{}, err
	}

	res, err := s.Service.Get(requisitionUUID)
	if err != nil {
		return uuid.UUID{}, err
	}

	if s.actor.Role() == auth.User && res.Author != s.actor.User() {
		return uuid.UUID{}, service.ErrUnauthorized
	}

	return s.Service.AddEntry(requisitionUUID, req)
}

func (s *Gate) RemoveEntry(requisitionUUID, entryUUID uuid.UUID) error {
	if err := service.Authorize(permLogged, s.actor); err != nil {
		return err
	}

	res, err := s.Service.Get(requisitionUUID)
	if err != nil {
		return err
	}

	if s.actor.Role() == auth.User && res.Author != s.actor.User() {
		return service.ErrUnauthorized
	}

	return s.Service.RemoveEntry(requisitionUUID, entryUUID)
}

func (s *Gate) Answer(requisitionUUID uuid.UUID, req requisition.AnswerRequisition) error {
	if err := service.Authorize(permAdmin, s.actor); err != nil {
		return err
	}

	return s.Service.Answer(requisitionUUID, req)
}

func (s *Gate) Cancel(requisitionUUID uuid.UUID) error {
	if err := service.Authorize(permLogged, s.actor); err != nil {
		return err
	}

	res, err := s.Service.Get(requisitionUUID)
	if err != nil {
		return err
	}

	if s.actor.Role() == auth.User && res.Author != s.actor.User() {
		return service.ErrUnauthorized
	}

	return s.Service.Cancel(requisitionUUID)
}

func (s *Gate) MarkFulfilled(requisitionUUID uuid.UUID) error {
	if err := service.Authorize(permPromoted, s.actor); err != nil {
		return err
	}

	return s.Service.MarkFulfilled(requisitionUUID)
}
