package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m-mizutani/octovy/pkg/service"
)

type Server struct{}

func New(svc *service.Service) *Server {
	return &Server{}
}

func (*Server) Listen(addr string, port int) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), r)

	return nil
}
