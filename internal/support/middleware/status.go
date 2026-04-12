package middleware

import "net/http"

type ResponseWriterWithStatus struct {
	http.ResponseWriter
	status int
}

func NewResponseWriterWithStatus(w http.ResponseWriter) *ResponseWriterWithStatus {
	return &ResponseWriterWithStatus{ResponseWriter: w, status: 200}
}

var _ http.ResponseWriter = (*ResponseWriterWithStatus)(nil)

func (rw *ResponseWriterWithStatus) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.status = status
}

func (rw *ResponseWriterWithStatus) StatusCode() int { return rw.status }

func (rw *ResponseWriterWithStatus) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
