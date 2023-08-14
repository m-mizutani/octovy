package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m-mizutani/octovy/pkg/usecase"
)

type Server struct {
	mux *chi.Mux
}

func New(uc *usecase.UseCase) *Server {
	r := chi.NewRouter()
	r.Use(preProcess)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})
	r.Route("/webhook", func(r chi.Router) {
		r.Post("/github", func(w http.ResponseWriter, r *http.Request) {

		})
	})

	return &Server{
		mux: r,
	}
}

func (x *Server) Mux() *chi.Mux {
	return x.mux
}
