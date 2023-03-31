// based on code taken from here (MIT license)
// https://github.com/go-chi/chi/blob/fd8a51eb979ab0fa87e6eb58b0e4f14428f31fd9/_examples/logging/main.go#L76

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

func NewStructuredLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{Logger: logger})
}

type StructuredLogger struct {
	Logger *slog.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	logger := l.Logger.With(
		slog.String("http_method", r.Method),
		slog.String("uri", fmt.Sprintf("%s%s", r.Host, r.RequestURI)),
	)

	entry := StructuredLoggerEntry{Logger: logger}

	return &entry
}

type StructuredLoggerEntry struct {
	Logger *slog.Logger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.Info("http request",
		slog.Int("resp_status", status),
	)
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger.Info("panic",
		slog.String("stack", string(stack)),
		slog.String("panic", fmt.Sprintf("%+v", v)),
	)
}
