package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"

	"time"

	log "github.com/seashell/drago/pkg/log"
)

// LoggingResponseWriter :
type LoggingResponseWriter struct {
	http.ResponseWriter
	code int
}

// NewLoggingResponseWriter :
func NewLoggingResponseWriter(rw http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{rw, http.StatusOK}
}

// WriteHeader :
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.code = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Hijack :
func (lrw *LoggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijacking not supported")
	}
	return h.Hijack()
}

// Logging :
func Logging(logger log.Logger) func(http.HandlerFunc) http.HandlerFunc {

	m := func(next http.HandlerFunc) http.HandlerFunc {

		return func(rw http.ResponseWriter, req *http.Request) {

			start := time.Now()
			url := req.RequestURI

			logger.Debugf("Request received (remote=%s, method=%s, path=%s)", req.RemoteAddr, req.Method, url)

			var status int
			defer func() {
				logger.Debugf("Request completed with status %d (method=%s, path=%s, duration=%s)", status, req.Method, url, time.Now().Sub(start))
			}()

			lrw := NewLoggingResponseWriter(rw)
			next(lrw, req)
			status = lrw.code
		}
	}

	return m
}
