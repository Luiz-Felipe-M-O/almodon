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
	"github.com/alan-b-lima/almodon/pkg/closer"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	items "github.com/alan-b-lima/almodon/internal/domain/item/resource"
	itemserve "github.com/alan-b-lima/almodon/internal/domain/item/service"
	itemstore "github.com/alan-b-lima/almodon/internal/domain/item/store"

	"github.com/alan-b-lima/almodon/internal/domain/material"
	materials "github.com/alan-b-lima/almodon/internal/domain/material/resource"
	materialserve "github.com/alan-b-lima/almodon/internal/domain/material/service"
	materialstore "github.com/alan-b-lima/almodon/internal/domain/material/store"

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

	"github.com/alan-b-lima/pkg/problem"
	"github.com/alan-b-lima/pkg/scheduler"

	_ "github.com/mattn/go-sqlite3"
)

type Almodon struct {
	http.ServeMux

	bundle closer.Bundle
}

type (
	Stores struct {
		Items      item.Store
		Materials  material.Store
		Promotions promotion.Store
		Sessions   session.Store
		Users      user.Store
	}

	Services struct {
		Auths      auth.Service
		Items      item.Service
		Materials  material.Service
		Promotions promotion.Service
		Sessions   session.Service
		Users      user.Service
	}

	Resources struct {
		Auth       *auths.Resource
		Items      *items.Resource
		Materials  *materials.Resource
		Promotions *promotions.Resource
		Users      *users.Resource
	}
)

var ErrNoRootUser = errors.New("no root user found in database")

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

	core := a.MountServices(stores)
	services := a.MountAuthServices(core)
	resources := a.MountResources(services)

	handlers := map[string]http.Handler{
		"auth":       resources.Auth,
		"items":      resources.Items,
		"materials":  resources.Materials,
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
		db, err = sql.Open("sqlite3", ".data/almodon.db")
	} else {
		db, err = sql.Open("sqlite3", "../.data/almodon.db")
	}
	if err != nil {
		return nil, err
	}
	a.bundle.Bundle(db)

	ctx := context.TODO()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	operations := [...]string{
		"PRAGMA foreign_keys = ON;",

		materialstore.Table,
		itemstore.Table,
		userstore.Table,
		sessionstore.Table,
		promotionstore.Table,

		materialstore.Indexes,
		promotionstore.Indexes,
		sessionstore.Indexes,
		userstore.Indexes,

		itemstore.Views,
		userstore.Views,
	}

	for _, op := range operations {
		if _, err := db.ExecContext(ctx, op); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (a *Almodon) MountSQLiteStores(db *sql.DB) (Stores, error) {
	stores := Stores{
		Items:      itemstore.New(db),
		Materials:  materialstore.New(db),
		Promotions: promotionstore.New(db),
		Sessions:   sessionstore.New(db),
		Users:      userstore.New(db),
	}

	if err := has_root_user(context.TODO(), stores.Users); err != nil {
		return Stores{}, err
	}

	now := time.Now()
	var (
		err_promotions = stores.Promotions.DeleteExpired(context.TODO(), now)
		err_sessions   = stores.Sessions.DeleteExpired(context.TODO(), now)
	)

	if err := problem.Join(err_sessions, err_promotions); err != nil {
		return Stores{}, err
	}

	return stores, nil
}

func (a *Almodon) MountServices(stores Stores) Services {
	scheduler := scheduler.New()
	a.bundle.BundleFunc(scheduler.Stop)
	scheduler.Start()

	services := Services{
		Items:     itemserve.New(stores.Items),
		Materials: materialserve.New(stores.Materials),
		Sessions:  sessionserve.New(stores.Sessions, scheduler),
		Users:     userserve.New(stores.Users),
	}

	services.Auths = authserve.New(services.Users, services.Sessions)
	services.Promotions = promotionserve.New(stores.Promotions, services.Users, scheduler)

	return services
}

func (a *Almodon) MountAuthServices(services Services) Services {
	authed := Services{
		Auths:      services.Auths,
		Items:      itemserve.NewGate(services.Items, services.Auths),
		Materials:  materialserve.NewGate(services.Materials, services.Auths),
		Promotions: promotionserve.NewGate(services.Promotions, services.Auths),
		Sessions:   services.Sessions,
		Users:      userserve.NewGate(services.Users, services.Auths),
	}

	return authed
}

func (a *Almodon) MountResources(services Services) Resources {
	resources := Resources{
		Auth:       auths.New(services.Auths),
		Items:      items.New(services.Items),
		Materials:  materials.New(services.Materials),
		Promotions: promotions.New(services.Promotions),
		Users:      users.New(services.Users),
	}

	return resources
}

func (a *Almodon) Close() error {
	return a.bundle.Close()
}

func has_root_user(ctx context.Context, store user.Store) error {
	_, err := store.GetBySIAPE(ctx, "0000000")
	if err != nil {
		if err == user.ErrNotFound {
			return ErrNoRootUser
		}

		return err
	}

	return nil
}
