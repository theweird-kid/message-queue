package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/theweird-kid/message-queue/internals/queue"
)

type Server struct {
	addr string
	r    *chi.Mux
	e    *queue.Exchange
}

func NewServer(addr string, exchange *queue.Exchange) *Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	return &Server{
		addr: addr,
		r:    router,
		e:    exchange,
	}
}

func (s *Server) HandleRoutes() {
	s.r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Home"))
	})

	s.r.Get("/topics", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("topics"))
	})
}

func (s *Server) Run() {
	log.Printf("Starting Server on port %v", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, s.r))
}
