package middleware

import (
	"net/http"
	"time"

	"github.com/marrria_mme/pr-reviewer-service/internal/transport/middleware/logctx"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger := logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"query":  r.URL.RawQuery,
		})

		ctx := logctx.WithLogger(r.Context(), logger)
		r = r.WithContext(ctx)

		logger.Info("request started")

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.WithFields(logrus.Fields{
			"status_code": rw.statusCode,
			"duration":    time.Since(start),
		}).Info("request completed")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
