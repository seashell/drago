package utils

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// realClock save request times
type realClock struct{}

func (rc *realClock) Now() time.Time {
	return time.Now()
}

func (rc *realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// LogOptions logging middleware options
type LogOptions struct {
}

// LoggingMiddleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type LoggingMiddleware struct {
	clock          timer
	logger         *log.Logger
	enableStarting bool
}

// NewLogger returns a new *LoggingMiddleware, yay!
func NewLogger(opts ...LogOptions) *LoggingMiddleware {

	logger := &log.Logger{
		Out:   os.Stderr,
		Level: log.DebugLevel,
		Formatter: &log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		},
	}

	return &LoggingMiddleware{
		clock:          &realClock{},
		logger:         logger,
		enableStarting: true,
	}
}

// realIP get the real IP from http request
func realIP(req *http.Request) string {
	ra := req.RemoteAddr
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		ra = strings.Split(ip, ", ")[0]
	} else if ip := req.Header.Get("X-Real-IP"); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

// Middleware implements mux middleware interface
func (m *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := m.clock.Now()

		reqID := r.Header.Get("X-Request-Id")
		remoteAddr := realIP(r)

		m.logger.Info(fmt.Sprintf("Request %s received from %s (%s %s)", reqID, remoteAddr, r.Method, r.RequestURI))

		lw := newLoggingResponseWriter(w)
		next.ServeHTTP(lw, r)

		latency := m.clock.Since(start)

		m.logger.Info(fmt.Sprintf("Request %s processed with status code %d (%s)", reqID, lw.statusCode, latency))
	})
}
