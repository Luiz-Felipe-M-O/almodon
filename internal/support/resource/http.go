package resource

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/alan-b-lima/pkg/problem"
)

func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if err, ok := errors.AsType[*problem.Error](err); ok {
		writeErrorJson(w, err, int(err.Kind))
		return
	}

	writeErrorJson(w, err, http.StatusInternalServerError)
}

func writeErrorJson(w http.ResponseWriter, err error, status int) {
	body, e := json.Marshal(err)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
}

var (
	reContentTypeApplicationJson = regexp.MustCompile(`^\s*(\*/\*|application/(json|\*))\s*(;.*)?\s*$`)
	reAcceptApplicationJson      = regexp.MustCompile(`(^|.*,)\s*(\*/\*|application/(json|\*))\s*(;.*)?\s*($|,.*)`)
)

func DecodeJSON(req any, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return ErrNoContentType
	}

	if !reContentTypeApplicationJson.MatchString(contentType) {
		return ErrUnsupportedContentType.Make(contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err == io.EOF {
			return ErrJSON.Message("unexpected end of input").Make()
		}

		return ErrJSON.Cause(err).Make()
	}

	return nil
}

var buffers = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func EncodeJSON(res any, status int, w http.ResponseWriter, r *http.Request) error {
	accept := r.Header.Get("Accept")
	if !reAcceptApplicationJson.MatchString(accept) {
		return ErrNotAcceptable.Make("application/json")
	}

	b := buffers.Get().(*bytes.Buffer)
	defer buffers.Put(b)
	b.Reset()

	if err := json.NewEncoder(b).Encode(res); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := io.Copy(w, b); err != nil {
		return err
	}

	return nil
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	WriteError(w, ErrResourceNotFound.Make(r.URL.Path))
}
