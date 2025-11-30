package requisitions

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/requisition"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Requisitions *auth.Gatekeeper[requisition.Service]

	Ident auth.Identifier
}

func New(requisitions requisition.Service) http.Handler {
	rc := Resource{
		Requisitions: auth.NewGatekeeper(requisitions),
		Ident:        requisitions,
	}

	routes := map[string]http.HandlerFunc{
		"GET /requisitions/{$}":                       rc.List,
		"GET /requisitions/{uuid}":                    rc.Get,
		"POST /requisitions/{$}":                      rc.Create,
		"PATCH /requisitions/{uuid}":                  rc.Patch,
		"DELETE /requisitions/{uuid}":                 rc.Delete,
		"POST /requisitions/{uuid}/entries/{$}":       rc.AddEntry,
		"DELETE /requisitions/{uuid}/entries/{entry}": rc.RemoveEntry,
		"POST /requisitions/{uuid}/answer/{$}":        rc.Answer,
		"POST /requisitions/{uuid}/cancel/{$}":        rc.Cancel,
		"POST /requisitions/{uuid}/fulfill/{$}":       rc.MarkFulfilled,
		"/":                                           resource.NotFound,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (requisition.ListResult, error) {
		req := requisition.ListParams{Offset: 0, Limit: 10}
		if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
			return requisition.ListResult{}, xerrors.ErrBadQueryParams.New(err)
		}

		ent, err := rc.Requisitions.Permit(act).List(req)
		if err != nil {
			return requisition.ListResult{}, err
		}

		res := requisition.ListResult{
			Offset:       ent.Offset,
			Length:       ent.Length,
			Records:      make([]requisition.Result, len(ent.Records)),
			TotalRecords: ent.TotalRecords,
		}
		for i := range len(ent.Records) {
			res.Records[i] = transform(&ent.Records[i])
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (requisition.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return requisition.Result{}, xerrors.ErrBadUUID
		}

		ent, err := rc.Requisitions.Permit(act).Get(uuid)
		if err != nil {
			return requisition.Result{}, err
		}

		return transform(&ent), nil
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(rc.Ident, func(act auth.Actor, req requisition.Create) (requisition.CreateResult, error) {
		res, err := rc.Requisitions.Permit(act).Create(req)
		if err != nil {
			return requisition.CreateResult{}, err
		}

		return requisition.CreateResult{UUID: res}, nil
	}, w, r)
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(rc.Ident, func(act auth.Actor, req requisition.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Requisitions.Permit(act).Patch(uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(rc.Ident, func(act auth.Actor) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Requisitions.Permit(act).Delete(uuid)
	}, w, r)
}

func (rc *Resource) AddEntry(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(rc.Ident, func(act auth.Actor, req requisition.AddEntry) (requisition.AddEntryResult, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return requisition.AddEntryResult{}, xerrors.ErrBadUUID
		}

		entryUUID, err := rc.Requisitions.Permit(act).AddEntry(uuid, req)
		if err != nil {
			return requisition.AddEntryResult{}, err
		}

		return requisition.AddEntryResult{UUID: entryUUID}, nil
	}, w, r)
}

func (rc *Resource) RemoveEntry(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(rc.Ident, func(act auth.Actor) error {
		reqUUID, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		entryUUID, err := uuid.FromString(r.PathValue("entry"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Requisitions.Permit(act).RemoveEntry(reqUUID, entryUUID)
	}, w, r)
}

func (rc *Resource) Answer(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(rc.Ident, func(act auth.Actor, req requisition.AnswerRequisition) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Requisitions.Permit(act).Answer(uuid, req)
	}, w, r)
}

func (rc *Resource) Cancel(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(rc.Ident, func(act auth.Actor, _ struct{}) (struct{}, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return struct{}{}, xerrors.ErrBadUUID
		}

		return struct{}{}, rc.Requisitions.Permit(act).Cancel(uuid)
	}, w, r)
}

func (rc *Resource) MarkFulfilled(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(rc.Ident, func(act auth.Actor, _ struct{}) (struct{}, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return struct{}{}, xerrors.ErrBadUUID
		}

		return struct{}{}, rc.Requisitions.Permit(act).MarkFulfilled(uuid)
	}, w, r)
}

func transform(e *requisition.Entity) requisition.Result {
	entries := make([]requisition.EntryResult, len(e.Entries))
	for i, entry := range e.Entries {
		entries[i] = requisition.EntryResult{
			UUID:     entry.UUID,
			Material: entry.Material,
			Quantity: entry.Quantity,
		}
	}

	answers := make([]requisition.AnswerResult, len(e.Answers))
	for i, answer := range e.Answers {
		answerEntries := make([]requisition.AnswerEntryResult, len(answer.Entries))
		for j, ae := range answer.Entries {
			answerEntries[j] = requisition.AnswerEntryResult{
				UUID:             ae.UUID,
				RequisitionEntry: ae.RequisitionEntry,
				ApprovedQuantity: ae.ApprovedQuantity,
				Notes:            ae.Notes,
			}
		}

		answers[i] = requisition.AnswerResult{
			UUID:       answer.UUID,
			Approver:   answer.Approver,
			Status:     answer.Status,
			Notes:      answer.Notes,
			Entries:    answerEntries,
			AnsweredAt: answer.AnsweredAt,
		}
	}

	return requisition.Result{
		UUID:        e.UUID,
		Author:      e.Author,
		Notes:       e.Notes,
		Destination: e.Destination,
		Status:      e.Status,
		Entries:     entries,
		Answers:     answers,
		Approver:    e.Approver,
		AnsweredAt:  e.AnsweredAt,
		Created:     e.Created,
		Updated:     e.Updated,
	}
}
