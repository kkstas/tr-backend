package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LogHTTP(logger *slog.Logger, next http.Handler) http.Handler {
	httpLogger := func(r *http.Request, lrw *loggingResponseWriter, start time.Time) {
		logger.LogAttrs(r.Context(),
			slog.LevelDebug,
			"request",
			slog.String("method", r.Method),
			slog.Int("status", lrw.statusCode),
			slog.String("uri", r.RequestURI),
			slog.String("duration", time.Since(start).String()),
		)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		now := time.Now()
		next.ServeHTTP(lrw, r)
		httpLogger(r, lrw, now)
	})
}
