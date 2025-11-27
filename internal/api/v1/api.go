package api

import (
	"errors"
	"net/http"

	itemrepo "github.com/alan-b-lima/almodon/internal/domain/item/repository"
	items "github.com/alan-b-lima/almodon/internal/domain/item/resource"
	itemserve "github.com/alan-b-lima/almodon/internal/domain/item/service"
	materialrepo "github.com/alan-b-lima/almodon/internal/domain/material/repository"
	materials "github.com/alan-b-lima/almodon/internal/domain/material/resource"
	materialserve "github.com/alan-b-lima/almodon/internal/domain/material/service"
	promotionrepo "github.com/alan-b-lima/almodon/internal/domain/promotion/repository"
	promotions "github.com/alan-b-lima/almodon/internal/domain/promotion/resource"
	promotionserve "github.com/alan-b-lima/almodon/internal/domain/promotion/service"
	sessionrepo "github.com/alan-b-lima/almodon/internal/domain/session/repository"
	sessionserve "github.com/alan-b-lima/almodon/internal/domain/session/service"
	userrepo "github.com/alan-b-lima/almodon/internal/domain/user/repository"
	users "github.com/alan-b-lima/almodon/internal/domain/user/resource"
	userserve "github.com/alan-b-lima/almodon/internal/domain/user/service"
	"github.com/alan-b-lima/almodon/pkg/closer"
)

type Handler struct {
	http.ServeMux
	cleanup closer.Bundle
}

func New() (*Handler, error) {
	var h Handler

	var (
		RepoItems, errRepoItems         = itemrepo.NewPersistentMap("../.data/items.json")
		RepoMaterials, errRepoMaterials = materialrepo.NewPersistentMap("../.data/materials.json")
		RepoPromotions                  = promotionrepo.NewMap()
		RepoSessions                    = sessionrepo.NewMap()
		RepoUsers, errRepoUsers         = userrepo.NewPersistantMap("../.data/users.json")
	)
	err := errors.Join(errRepoItems, errRepoMaterials, errRepoUsers)
	if err != nil {
		closer.CloseMany(RepoItems, RepoMaterials, RepoPromotions, RepoSessions, RepoUsers)
		return nil, err
	}

	var (
		CoreItems      = &itemserve.Core{RepoItems}
		CoreMaterials  = &materialserve.Core{RepoMaterials}
		CorePromotions = &promotionserve.Core{RepoPromotions, nil}
		CoreSessions   = &sessionserve.Core{RepoSessions}
		CoreUsers      = &userserve.Core{RepoUsers, CoreSessions, CorePromotions}
	)
	CorePromotions.Users = CoreUsers

	var (
		ServiceItems      = itemserve.New(CoreItems)
		ServiceMaterials  = materialserve.New(CoreMaterials)
		ServicePromotions = promotionserve.New(CorePromotions)
		ServiceUsers      = userserve.New(CoreUsers)
	)

	var (
		items      = items.New(ServiceItems, ServiceUsers)
		materials  = materials.New(ServiceMaterials, ServiceUsers)
		promotions = promotions.New(ServicePromotions, ServiceUsers)
		users      = users.New(ServiceUsers)
	)

	resources := map[string]http.Handler{
		"items":      items,
		"materials":  materials,
		"promotions": promotions,
		"users":      users,
	}

	for name, handler := range resources {
		h.Handle("/api/v1/"+name+"/", http.StripPrefix("/api/v1", handler))
	}

	h.cleanup.BundleMany(
		RepoItems, RepoMaterials, RepoPromotions, RepoSessions, RepoUsers,
		CoreItems, CoreMaterials, CorePromotions, CoreSessions, CoreUsers,
		ServiceItems, ServiceMaterials, ServicePromotions, ServiceUsers,
		items, materials, promotions, users,
	)

	return &h, nil
}

func (h *Handler) Close() error {
	return h.cleanup.Close()
}
