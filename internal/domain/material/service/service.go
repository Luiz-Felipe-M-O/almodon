package materialserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service struct {
	materials material.Repository
}

func NewService(materials material.Repository) material.Service {
	return &Service{
		materials: materials,
	}
}

func (s *Service) List(act auth.Actor, req material.ListRequest) (material.ListResponse, error) {
	res, err := material.List(s.materials, req.Offset, req.Limit)
	if err != nil {
		return material.ListResponse{}, err
	}

	lres := material.ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]material.Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := range res.Records {
		transformP(&lres.Records[i], &res.Records[i])
	}

	return lres, nil
}

func (s *Service) Get(act auth.Actor, req material.GetRequest) (material.Response, error) {
	res, err := material.Get(s.materials, req.UUID)
	if err != nil {
		return material.Response{}, err
	}

	return transform(&res), nil
}

func (s *Service) ListBySIADS(act auth.Actor, req material.ListBySIADSRequest) (material.ListResponse, error) {
	res, err := material.ListBySIADS(s.materials, req.SIADS)
	if err != nil {
		return material.ListResponse{}, err
	}

	lres := material.ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]material.Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := range res.Records {
		transformP(&lres.Records[i], &res.Records[i])
	}

	return lres, nil
}

func (s *Service) ListByCATMAT(act auth.Actor, req material.ListByCATMATRequest) (material.ListResponse, error) {
	res, err := material.ListByCATMAT(s.materials, req.CATMAT)
	if err != nil {
		return material.ListResponse{}, err
	}

	lres := material.ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]material.Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := range res.Records {
		transformP(&lres.Records[i], &res.Records[i])
	}

	return lres, nil
}

func (s *Service) ListByECampus(act auth.Actor, req material.ListByECampusRequest) (material.ListResponse, error) {
	res, err := material.ListByECAMPUS(s.materials, req.ECAMPUS)
	if err != nil {
		return material.ListResponse{}, err
	}

	lres := material.ListResponse{
		Offset:       res.Offset,
		Length:       res.Length,
		Records:      make([]material.Response, res.Length),
		TotalRecords: res.TotalRecords,
	}
	for i := range res.Records {
		transformP(&lres.Records[i], &res.Records[i])
	}

	return lres, nil
}

func (s *Service) Create(act auth.Actor, req material.CreateRequest) (uuid.UUID, error) {
	return material.Create(
		s.materials,
		req.Name,
		req.SIADS,
		req.CATMAT,
		req.ECAMPUS,
		req.Description,
		req.Unit,
		req.MinQuantity,
	)
}

func (s *Service) Patch(act auth.Actor, req material.PatchRequest) error {
	var float64Opt opt.Opt[float64]

	return material.Patch(
		s.materials,
		req.UUID,
		req.Name,
		req.SIADS,
		req.CATMAT,
		req.ECAMPUS,
		req.Description,
		req.Unit,
		float64Opt,
	)
}

func (s *Service) UpdateMinQuantity(act auth.Actor, req material.UpdateMinQuantityRequest) error {
	var stringOpt opt.Opt[string]

	return material.Patch(
		s.materials,
		req.UUID,
		stringOpt,
		stringOpt,
		stringOpt,
		stringOpt,
		stringOpt,
		stringOpt,
		opt.Some(req.MinQuantity),
	)
}

func (s *Service) Delete(act auth.Actor, req material.DeleteRequest) error {
	return material.Delete(s.materials, req.UUID)
}

func transform(e *material.Entity) material.Response {
	return material.Response{
		UUID:        e.UUID,
		Name:        e.Name,
		SIADS:       e.SIADS,
		CATMAT:      e.CATMAT,
		ECAMPUS:     e.ECAMPUS,
		Description: e.Description,
		Unit:        e.Unit,
		MinQuantity: e.MinQuantity,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func transformP(r *material.Response, e *material.Entity) {
	r.UUID = e.UUID
	r.Name = e.Name
	r.SIADS = e.SIADS
	r.CATMAT = e.CATMAT
	r.ECAMPUS = e.ECAMPUS
	r.Description = e.Description
	r.Unit = e.Unit
	r.MinQuantity = e.MinQuantity
	r.CreatedAt = e.CreatedAt
	r.UpdatedAt = e.UpdatedAt
}
