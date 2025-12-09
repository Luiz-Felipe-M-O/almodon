package auth

import "github.com/alan-b-lima/almodon/pkg/uuid"

type Identifier interface {
	Actor(uuid.UUID) (Actor, error)
}

type Gatekeeper[S any] struct {
	Service S
}

func NewGatekeeper[S any](service S) *Gatekeeper[S] {
	return &Gatekeeper[S]{
		Service: service,
	}
}

func (g *Gatekeeper[S]) Permit(act Actor) S {
	var service any = g.Service

	if allower, ok := service.(interface{ Allow(Actor) S }); ok {
		return allower.Allow(act)
	}

	return g.Service
}
