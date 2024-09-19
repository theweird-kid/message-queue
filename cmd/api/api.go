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

type Handler struct {
	e *queue.Exchange
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
	h := &Handler{
		e: s.e,
	}
	s.r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Home"))
	})

	// Get the list of available topics
	s.r.Get("/topics", h.GetTopics)
	// Create a New Topic
	s.r.Post("/topic", h.CreateTopic)
	// Publish message to a Topic
	s.r.Post("/pub", h.PublishMessage)

}

func (s *Server) Run() {
	log.Printf("Starting Server on port %v", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, s.r))
}
