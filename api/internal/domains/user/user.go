package user

import (
	"github.com/Corray333/keep_it/internal/domains/user/storage"
	"github.com/Corray333/keep_it/internal/domains/user/transport"
	"github.com/Corray333/keep_it/internal/global_storage"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	store global_storage.Storage
}

func NewController() *Controller {
	return &Controller{}
}

func (c Controller) Init(router *chi.Mux, storeGlobal global_storage.Storage) {

	store := storage.NewStorage(storeGlobal.GetDB(), storeGlobal.GetRedis())
	router.Post("/api/users/login", transport.LogIn(store))
	router.Post("/api/users/signup", transport.SignUp(store))
	router.Get("/api/users/refresh", transport.RefreshAccessToken(store))
	router.Get("/api/users/check-username", transport.CheckUsername(store))
	router.Post("/api/users/check-code", transport.CheckCode(store))
	router.With(auth.NewAuthMiddleware()).Get("/api/users/{id}", transport.GetUser(store))
}
