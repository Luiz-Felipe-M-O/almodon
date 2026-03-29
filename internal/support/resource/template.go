package resource

import (
	"context"
	"net/http"
)

func GetHandler[TRes any](ctx context.Context, proc func(context.Context) (TRes, error), w http.ResponseWriter, r *http.Request) {
	ctx, err := Session(ctx, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	res, err := proc(ctx)
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		WriteError(w, err)
		return
	}
}

func PostHandler[TRes, TReq any](ctx context.Context, proc func(context.Context, TReq) (TRes, error), w http.ResponseWriter, r *http.Request) {
	ctx, err := Session(ctx, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	var req TReq
	if err := DecodeJSON(&req, r); err != nil {
		WriteError(w, err)
		return
	}

	uuid, err := proc(ctx, req)
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := EncodeJSON(&uuid, http.StatusCreated, w, r); err != nil {
		WriteError(w, err)
		return
	}
}

func PutHandler[TReq any](ctx context.Context, proc func(context.Context, TReq) error, w http.ResponseWriter, r *http.Request) {
	ctx, err := Session(ctx, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	var req TReq
	if err := DecodeJSON(&req, r); err != nil {
		WriteError(w, err)
		return
	}

	if err := proc(ctx, req); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteHandler(ctx context.Context, proc func(context.Context) error, w http.ResponseWriter, r *http.Request) {
	ctx, err := Session(ctx, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := proc(ctx); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
