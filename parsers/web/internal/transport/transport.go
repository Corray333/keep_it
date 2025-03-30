package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/parsers/web/pkg/server/auth"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

type Transport struct {
	service service
	router  *chi.Mux
}

type service interface {
	ProcessHTML(ctx context.Context, userID int64, document string, url string) error
}

func newRouter() *chi.Mux {
	router := chi.NewMux()
	// router.Use(logger.NewLoggerMiddleware())

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

	return router
}

func New(service service) *Transport {
	router := newRouter()

	t := &Transport{
		router:  router,
		service: service,
	}

	t.RegisterRoutes()

	return t
}

func (t *Transport) Run() {
	slog.Info("Starting server on port " + viper.GetString("server.port") + "...")
	if err := http.ListenAndServe("0.0.0.0"+":"+viper.GetString("server.port"), t.router); err != nil {
		slog.Error("Failed to start server", "error", err)
		panic(err)
	}
}

func (t *Transport) RegisterRoutes() {

	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())
	})
	t.router.Post("/api/parser/web", t.processWeb)
}

type processWebRequest struct {
	Document string `json:"document"`
	Url      string `json:"url"`
}

func (t *Transport) processWeb(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test")
	var req processWebRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(req.Url)
	if err != nil {
		http.Error(w, "failed to fetch URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	doc, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response body", http.StatusInternalServerError)
		return
	}
	req.Document = string(doc)

	// userID, ok := r.Context().Value(auth.CtxUserIDKey).(int64)
	// if !ok {
	// 	slog.Error("user id not found in context")
	// 	http.Error(w, "user id not found in context", http.StatusInternalServerError)
	// 	return
	// }

	err = t.service.ProcessHTML(r.Context(), 0, req.Document, req.Url)
	if err != nil {
		http.Error(w, "failed to process HTML", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
