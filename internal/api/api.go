package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	auths "github.com/alan-b-lima/almodon/internal/domain/auth/resource"
	authserve "github.com/alan-b-lima/almodon/internal/domain/auth/service"

	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	promotions "github.com/alan-b-lima/almodon/internal/domain/promotion/resource"
	promotionserve "github.com/alan-b-lima/almodon/internal/domain/promotion/service"
	promotionstore "github.com/alan-b-lima/almodon/internal/domain/promotion/store"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	sessionserve "github.com/alan-b-lima/almodon/internal/domain/session/service"
	sessionstore "github.com/alan-b-lima/almodon/internal/domain/session/store"

	"github.com/alan-b-lima/almodon/internal/domain/user"
	users "github.com/alan-b-lima/almodon/internal/domain/user/resource"
	userserve "github.com/alan-b-lima/almodon/internal/domain/user/service"
	userstore "github.com/alan-b-lima/almodon/internal/domain/user/store"

	"github.com/alan-b-lima/almodon/internal/support/store"

	"github.com/alan-b-lima/almodon/pkg/closer"

	"github.com/alan-b-lima/pkg/problem"
	"github.com/alan-b-lima/pkg/scheduler"

	_ "modernc.org/sqlite"
)

type Almodon struct {
	http.ServeMux

	cleanup closer.Bundle
}

type (
	Stores struct {
		Promotions promotion.Store
		Sessions   session.Store
		Users      user.Store
	}

	Services struct {
		Auths      auth.Service
		Promotions promotion.Service
		Sessions   session.Service
		Users      user.Service
	}

	Resources struct {
		Auth       *auths.Resource
		Promotions *promotions.Resource
		Users      *users.Resource
	}
)

func New() (*Almodon, error) {
	var a Almodon
	var err error

	defer func(err *error) {
		if *err != nil {
			a.Close()
		}
	}(&err)

	db, err := a.MountSQLiteDB()
	if err != nil {
		return nil, err
	}

	stores, err := a.MountSQLiteStores(db)
	if err != nil {
		return nil, err
	}

	services, err := a.MountServices(stores)
	if err != nil {
		return nil, err
	}

	resources := a.MountResources(services)

	handlers := map[string]http.Handler{
		"auth":       resources.Auth,
		"promotions": resources.Promotions,
		"users":      resources.Users,
	}
	for name, handler := range handlers {
		a.Handle("/api/v1/"+name+"/", http.StripPrefix("/api/v1", handler))
	}

	return &a, nil
}

func (a *Almodon) MountSQLiteDB() (*sql.DB, error) {
	var db *sql.DB

	_, err := os.Stat(".data/almodon.db")
	if !errors.Is(err, os.ErrNotExist) {
		db, err = sql.Open("sqlite", ".data/almodon.db")
	} else {
		db, err = sql.Open("sqlite", "../.data/almodon.db")
	}
	if err != nil {
		return nil, err
	}
	a.cleanup.Bundle(db)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	ctx := context.TODO()
	err = store.WithTx(ctx, db, func(tx store.DBTx) error {
		operations := [...]string{
			userstore.Table,
			sessionstore.Table,
			promotionstore.Table,

			userstore.Indexes,
			sessionstore.Indexes,
			promotionstore.Indexes,

			userstore.Views,
		}

		for _, op := range operations {
			if _, err := tx.ExecContext(ctx, op); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (a *Almodon) MountSQLiteStores(db *sql.DB) (Stores, error) {
	stores := Stores{
		Promotions: promotionstore.New(db),
		Sessions:   sessionstore.New(db),
		Users:      userstore.New(db),
	}
	a.cleanup.BundleMany(
		stores.Promotions,
		stores.Sessions,
		stores.Users,
	)

	now := time.Now()
	var (
		err_promotions = stores.Promotions.DeleteExpired(context.Background(), now)
		err_sessions   = stores.Sessions.DeleteExpired(context.Background(), now)
	)

	err := problem.Join(err_sessions, err_promotions)
	if err != nil {
		return Stores{}, err
	}

	return stores, nil
}

func (a *Almodon) MountServices(stores Stores) (Services, error) {
	scheduler := scheduler.New()
	a.cleanup.BundleFunc(scheduler.Stop)
	scheduler.Start()

	var (
		sessions, err_sessions     = sessionserve.New(stores.Sessions, scheduler)
		users                      = userserve.New(stores.Users)
		promotions, err_promotions = promotionserve.New(stores.Promotions, users, scheduler)
		auths                      = authserve.New(users, sessions)
	)
	a.cleanup.BundleMany(auths, promotions, sessions, users)

	err := problem.Join(err_sessions, err_promotions)
	if err != nil {
		return Services{}, err
	}

	services := Services{
		Auths:      auths,
		Promotions: promotionserve.NewGate(promotions, auths),
		Sessions:   sessions,
		Users:      userserve.NewGate(users, auths),
	}
	a.cleanup.BundleMany(
		services.Auths,
		services.Promotions,
		services.Sessions,
		services.Users,
	)

	return services, nil
}

func (a *Almodon) MountResources(services Services) Resources {
	resources := Resources{
		Auth:       auths.New(services.Auths),
		Promotions: promotions.New(services.Promotions),
		Users:      users.New(services.Users),
	}
	a.cleanup.BundleMany(
		resources.Auth,
		resources.Promotions,
		resources.Users,
	)

	return resources
}

func (a *Almodon) Close() error {
	return a.cleanup.Close()
}
