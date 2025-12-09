package resource

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
)

func GetHandler[TRes any](gk auth.Identifier, proc func(auth.Actor) (TRes, error), w http.ResponseWriter, r *http.Request) {
	act, err := Session(gk, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	res, err := proc(act)
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		WriteError(w, err)
		return
	}
}

func PostHandler[TRes, TReq any](gk auth.Identifier, proc func(act auth.Actor, req TReq) (TRes, error), w http.ResponseWriter, r *http.Request) {
	act, err := Session(gk, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	var req TReq
	if err := DecodeJSON(&req, r); err != nil {
		WriteError(w, err)
		return
	}

	uuid, err := proc(act, req)
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := EncodeJSON(&uuid, http.StatusCreated, w, r); err != nil {
		WriteError(w, err)
		return
	}
}

func PutHandler[TReq any](gk auth.Identifier, proc func(act auth.Actor, req TReq) error, w http.ResponseWriter, r *http.Request) {
	act, err := Session(gk, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	var req TReq
	if err := DecodeJSON(&req, r); err != nil {
		WriteError(w, err)
		return
	}

	if err := proc(act, req); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteHandler(auth auth.Identifier, proc func(act auth.Actor) error, w http.ResponseWriter, r *http.Request) {
	act, err := Session(auth, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := proc(act); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
