package api

import (
	"net/http"

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
		RepoPromotions          = promotionrepo.NewMap()
		RepoSessions            = sessionrepo.NewMap()
		RepoUsers, errRepoUsers = userrepo.NewPersistantMap("../.data/users.json")
	)
	if errRepoUsers != nil {
		return nil, errRepoUsers
	}

	var (
		CorePromotions = &promotionserve.Core{RepoPromotions, nil}
		CoreSessions   = &sessionserve.Core{RepoSessions}
		CoreUsers      = &userserve.Core{RepoUsers, CoreSessions, CorePromotions}
	)
	CorePromotions.Users = CoreUsers

	var (
		ServicePromotions = promotionserve.New(CorePromotions)
		ServiceUsers      = userserve.New(CoreUsers)
	)

	var (
		promotions = promotions.New(ServicePromotions, ServiceUsers)
		users      = users.New(ServiceUsers)
	)

	resources := map[string]http.Handler{
		"promotions": promotions,
		"users":      users,
	}

	for name, handler := range resources {
		h.Handle("/api/v1/"+name+"/", http.StripPrefix("/api/v1", handler))
	}

	h.cleanup.BundleMany(
		RepoPromotions,
		RepoSessions,
		errRepoUsers,
		CorePromotions,
		CoreSessions,
		CoreUsers,
		ServicePromotions,
		ServiceUsers,
		promotions,
		users,
	)

	return &h, nil
}

func (h *Handler) Close() error {
	return h.cleanup.Close()
}
