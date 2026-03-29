package sse

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

// ServerSentEvent is a struct that represents a Server-Sent Event connection.
// It provides methods to safely write data to the connection and dispatch
// events to the client.
//
// ServerSentEvent's are NOT safe for use by multiple goroutines.
//
// See [https://html.spec.whatwg.org/multipage/server-sent-events.html] for
// specification details.
type ServerSentEvent struct {
	conn interface {
		io.Writer
		http.Flusher
	}

	id  string
	typ string
	buf bytes.Buffer

	idsent  bool
	typsent bool
	bufsent bool
}

var ErrUnsupported = errors.New("sse: streaming unsupported")

// New creates a new ServerSentEvent and attaches it to the provided
// [http.ResponseWriter].
//
// It returns an error if the provided http.ResponseWriter does not support
// the necessary interfaces for streaming.
func New(w http.ResponseWriter) (*ServerSentEvent, error) {
	var e ServerSentEvent
	if err := e.Attach(w); err != nil {
		return nil, err
	}

	return &e, nil
}

// Attach attaches the ServerSentEvent to the provided [http.ResponseWriter].
//
// It returns an error if the provided http.ResponseWriter does not support
// the necessary interfaces for streaming.
func (e *ServerSentEvent) Attach(w http.ResponseWriter) error {
	wf, ok := w.(interface {
		http.ResponseWriter
		http.Flusher
	})
	if !ok {
		return ErrUnsupported
	}

	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")

	w.WriteHeader(http.StatusOK)

	e.conn = wf

	e.id = w.Header().Get("Last-Event-Id")
	e.typ = ""

	e.idsent = true
	e.typsent = true
	e.bufsent = true

	return nil
}

// Detach detaches the ServerSentEvent from its current connection.
//
// Detaching will be reset the ServerSentEvent to its zero value, with no ID,
// type, or buffered data.
//
// Detaching does not close the underlying connection, but it will prevent any
// further dispatches from affecting the connection.
//
// Detach is not required to be run before attaching to a new connection, but
// may be used to allow the underlying response writer to be claimed by the
// garbage collector.
func (e *ServerSentEvent) Detach() {
	*e = ServerSentEvent{}
}

// Type returns the current event type of the ServerSentEvent, or an empty string
// if no type has been set.
func (e *ServerSentEvent) Type() string { return e.typ }

// ID returns the current event ID of the ServerSentEvent, or an empty string if
// no ID has been set.
func (e *ServerSentEvent) ID() string { return e.id }

// SetType sets the event type of the ServerSentEvent.
//
// The event type will be sent to the client on the next dispatch if it has
// been set, is not empty and not already sent.
//
// If the provided type contains any of line feed (LF U+000A) or carriage
// return (CR U+000D), only the substring up to the first these, exclusive,
// will be used, and the rest will be ignored.
func (e *ServerSentEvent) SetType(typ string) {
	index := strings.IndexAny(typ, "\r\n")
	if index == -1 {
		index = len(id)
	}

	e.typsent = false
	e.typ = typ[:index]
}

// SetID sets the event ID of the ServerSentEvent.
//
// The event ID will be sent to the client on the next dispatch if it has been
// set and not already sent.
func (e *ServerSentEvent) SetID(id string) {
	index := strings.IndexAny(id, "\r\n")
	if index == -1 {
		index = len(id)
	}

	e.idsent = false
	e.id = id[:index]
}

var (
	id    = []byte("id: ")
	event = []byte("event: ")
	data  = []byte("data: ")

	idping   = []byte("id\r\n")
	dataping = []byte("data\r\n\r\n")

	crlf     = []byte("\r\n")
	crlfcrlf = []byte("\r\n\r\n")
)

// Write writes the provided byte slice to the ServerSentEvent's buffer,
// properly formatting it according to the Server-Sent Events specification.
//
// Server-Sent Events SHOULD NOT be used to send arbitrary binary data, the
// stream is always interpreted as UTF-8; and carriage returns U+000D are
// impossible to send, therefore ignored here to avoid problems down the line.
//
// The provided byte slice will be automatically split into multiple data lines
// if it contains newline characters.
//
// err is always nil. If the buffer becomes too large, Write will panic with
// [bytes.ErrTooLarge].
//
// Write does not automatically dispatch the event to the client, you must call
// Dispatch() to send the buffered data.
func (e *ServerSentEvent) Write(b []byte) (int, error) {
	e.bufsent = false
	if len(b) == 0 {
		return 0, nil
	}

	e.buf.Grow(len(b))
	orglen := e.buf.Len()

	if e.buf.Len() == 0 {
		e.buf.Write(data)
	}

	for i := 0; i < len(b); i++ {
		switch b[i] {
		case '\r':
			e.buf.Write(b[:i])
			b = b[i+1:]
			i = -1

		case '\n':
			e.buf.Write(b[:i])
			e.buf.Write(crlf)
			e.buf.Write(data)

			b = b[i+1:]
			i = -1
		}
	}

	e.buf.Write(b)
	return orglen - e.buf.Len(), nil
}

// Dispatch sends the buffered event data to the client, properly formatting it
// according to the Server-Sent Events specification.
//
// Dispatch will also send the event's ID and type if they have been set and
// not already sent.
//
// After dispatching, the ServerSentEvent's buffer and type will be cleared,
// but the ID will remain set.
func (e *ServerSentEvent) Dispatch() (int, error) {
	n, err := e.dispatch()
	if n > 0 {
		e.conn.Flush()
	}

	return n, err
}

func (e *ServerSentEvent) dispatch() (int, error) {
	var ntotal int

	if !e.idsent {
		var err error
		if e.id != "" {
			err = write(e.conn, &ntotal, id, []byte(e.id), crlf)
		} else {
			err = write(e.conn, &ntotal, idping)
		}

		if err != nil {
			return ntotal, err
		}

		e.idsent = true
	}

	if !e.typsent && e.typ != "" {
		err := write(e.conn, &ntotal, event, []byte(e.typ), crlf)
		if err != nil {
			return ntotal, err
		}

		e.typ = ""
		e.typsent = true
	}

	if !e.bufsent {
		if e.buf.Len() == 0 {
			err := write(e.conn, &ntotal, dataping)
			return ntotal, err
		}

		err := write(e.conn, &ntotal, e.buf.Bytes(), crlfcrlf)
		if err != nil {
			return ntotal, err
		}

		e.buf.Reset()
		e.bufsent = true
	}

	return ntotal, nil
}

var comment = []byte(":")

// Comment writes the provided byte slice as a comment to the ServerSentEvent's
// connection, properly formatting it according to the Server-Sent Events
// specification.
//
// Comments have no special meaning and are ignored by the client, but they can
// be used to keep the connection alive. You may use Comment(nil) to send a
// ping comment, which is a comment with no content.
//
// The same said about arbitrary binary data in Write() applies to Comment().
func (e *ServerSentEvent) Comment(b []byte) (int, error) {
	defer e.conn.Flush()
	var ntotal int

	err := write(e.conn, &ntotal, comment)
	if err != nil {
		return ntotal, err
	}

	if len(b) == 0 {
		return ntotal, nil
	}

	for i := 0; i < len(b); i++ {
		switch b[i] {
		case '\r':
			err := write(e.conn, &ntotal, b[:i])
			if err != nil {
				return ntotal, err
			}

			b = b[i+1:]
			i = -1

		case '\n':
			err := write(e.conn, &ntotal, b[:i], crlf, comment)
			if err != nil {
				return ntotal, err
			}

			b = b[i+1:]
			i = -1
		}
	}

	err = write(e.conn, &ntotal, b)
	return ntotal, err
}

func write(w io.Writer, n *int, bufs ...[]byte) error {
	for _, b := range bufs {
		m, err := w.Write(b)
		*n += m

		if err != nil {
			return err
		}
	}

	return nil
}
