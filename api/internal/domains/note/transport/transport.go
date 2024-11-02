package transport

import (
	"github.com/go-chi/chi/v5"
)

type NoteTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
}

func New(router *chi.Mux, service service) *NoteTransport {
	return &NoteTransport{
		router:  router,
		service: service,
	}
}

func (t *NoteTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {

	})
}
