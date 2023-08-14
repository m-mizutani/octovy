package server

import (
	"net/http"

	"log/slog"

	"github.com/google/uuid"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func preProcess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := utils.Logger().With(slog.String("request_id", uuid.NewString()))

		ctx := model.NewContext(
			model.WithLogger(logger),
			model.WithBase(r.Context()),
		)

		lw := &statusCodeLogger{ResponseWriter: w}
		next.ServeHTTP(lw, r.WithContext(ctx))

		logger.Info("http access",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Int("status_code", lw.statusCode),
			slog.Int64("content_length", r.ContentLength),
			slog.String("user_agent", r.UserAgent()),
			slog.String("referer", r.Referer()),
		)
	})
}

type statusCodeLogger struct {
	http.ResponseWriter
	statusCode int
}

func (x *statusCodeLogger) WriteHeader(code int) {
	x.statusCode = code
	x.ResponseWriter.WriteHeader(code)
}
