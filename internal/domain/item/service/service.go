package itemserve

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service struct {
	items item.Repository
}

func NewService(items item.Repository) item.Service {
	return &Service{
		items: items,
	}
}

func (s *Service) List(act auth.Actor, req item.ListRequest) (item.ListResponse, error) {
	res, err := item.List(s.items, req.Offset, req.Limit)
	if err != nil {
		return item.ListResponse{}, err
	}

	lres := item.ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]item.Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := range res.Records {
		transformP(&lres.Records[i], &res.Records[i])
	}

	return lres, nil
}

func (s *Service) Get(act auth.Actor, req item.GetRequest) (item.Response, error) {
	res, err := item.Get(s.items, req.UUID)
	if err != nil {
		return item.Response{}, err
	}

	return transform(&res), nil
}

func (s *Service) GetByBatch(act auth.Actor, req item.GetByBatchRequest) (item.ListResponse, error) {
	res, err := item.ListByBatch(s.items, req.Batch)
	if err != nil {
		return item.ListResponse{}, err
	}

	lres := item.ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]item.Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := range res.Records {
		transformP(&lres.Records[i], &res.Records[i])
	}

	return lres, nil
}

func (s *Service) GetByMaterial(act auth.Actor, req item.GetByMaterialRequest) (item.ListResponse, error) {
	res, err := item.ListByMaterial(s.items, req.Material)
	if err != nil {
		return item.ListResponse{}, err
	}

	lres := item.ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]item.Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := range res.Records {
		transformP(&lres.Records[i], &res.Records[i])
	}

	return lres, nil
}

func (s *Service) Create(act auth.Actor, req item.CreateRequest) (uuid.UUID, error) {
	return item.Create(
		s.items,
		req.Batch,
		req.Material,
		req.Quantity,
		req.Expiration,
	)
}

func (s *Service) Patch(act auth.Actor, req item.PatchRequest) error {
	var float64Opt opt.Opt[float64]
	var timeOpt opt.Opt[time.Time]

	return item.Patch(
		s.items,
		req.UUID,
		req.Batch,
		req.Material,
		float64Opt,
		timeOpt,
	)
}

func (s *Service) UpdateQuantity(act auth.Actor, req item.UpdateQuantityRequest) error {
	var uuidOpt opt.Opt[uuid.UUID]
	var timeOpt opt.Opt[time.Time]

	return item.Patch(
		s.items,
		req.UUID,
		uuidOpt,
		uuidOpt,
		opt.Some(req.Quantity),
		timeOpt,
	)
}

func (s *Service) Delete(act auth.Actor, req item.DeleteRequest) error {
	return item.Delete(s.items, req.UUID)
}

func transform(e *item.Entity) item.Response {
	isExpired := false
	if !e.Expiration.IsZero() && time.Now().After(e.Expiration) {
		isExpired = true
	}

	return item.Response{
		UUID:       e.UUID,
		Batch:      e.Batch,
		Material:   e.Material,
		Quantity:   e.Quantity,
		Expiration: e.Expiration,
		IsExpired:  isExpired,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func transformP(r *item.Response, e *item.Entity) {
	isExpired := false
	if !e.Expiration.IsZero() && time.Now().After(e.Expiration) {
		isExpired = true
	}

	r.UUID = e.UUID
	r.Batch = e.Batch
	r.Material = e.Material
	r.Quantity = e.Quantity
	r.Expiration = e.Expiration
	r.IsExpired = isExpired
	r.CreatedAt = e.CreatedAt
	r.UpdatedAt = e.UpdatedAt
}
