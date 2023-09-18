package server

import (
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/m-mizutani/octovy/pkg/utils"
)

type Server struct {
	mux *chi.Mux
}

func safeWrite(w http.ResponseWriter, code int, body []byte) {
	w.WriteHeader(code)
	if _, err := w.Write(body); err != nil {
		utils.Logger().Error("fail to write response", slog.Any("error", err))
	}
}

func New(uc usecase.UseCase, secret types.GitHubAppSecret) *Server {
	r := chi.NewRouter()
	r.Use(preProcess)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		safeWrite(w, http.StatusOK, []byte("ok"))
	})
	r.Route("/webhook", func(r chi.Router) {
		r.Post("/github", func(w http.ResponseWriter, r *http.Request) {
			if err := handleGitHubEvent(uc, r, secret); err != nil {
				utils.Logger().Warn("fail to handle GitHub event", slog.Any("error", err))
				safeWrite(w, http.StatusInternalServerError, []byte(err.Error()))
				return
			}

			safeWrite(w, http.StatusOK, []byte("ok"))
		})
	})

	return &Server{
		mux: r,
	}
}

func (x *Server) Mux() *chi.Mux {
	return x.mux
}
