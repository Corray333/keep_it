package transport

import (
	"github.com/go-chi/chi/v5"
)

type UserTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
}

func New(router *chi.Mux, service service) *UserTransport {
	return &UserTransport{
		router:  router,
		service: service,
	}
}

func (t *UserTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {

	})
}
