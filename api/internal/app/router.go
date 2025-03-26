package app

import (
	"github.com/Corray333/keep_it/pkg/server/logger"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func newRouter() *chi.Mux {
	router := chi.NewMux()
	router.Use(logger.NewLoggerMiddleware())

	// TODO: get allowed origins, headers and methods from cfg
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Set-Cookie", "Refresh", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(c.Handler)

	router.Get("/api/swagger/*", httpSwagger.WrapHandler)

	return router
}
