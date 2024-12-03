package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/category/entities"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type CategoryTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
	GetCategories(ctx context.Context, userID int64) ([]entities.Category, error)
}

func New(router *chi.Mux, service service) *CategoryTransport {

	return &CategoryTransport{
		router:  router,
		service: service,
	}
}

func (t *CategoryTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())

		r.Post("/api/categories", t.createCategory)
		r.Get("/api/categories", t.getCategories)
	})
}

func (t *CategoryTransport) Run() {

}

func (t *CategoryTransport) createCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		slog.Error("User id not found in context")
		helpers.SendError(w, helpers.ErrInternal)
		return
	}

	category := &entities.Category{}
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		helpers.SendError(w, err)
		return
	}

	category.OwnerID = userID

	category, err := t.service.CreateCategory(ctx, category)
	if err != nil {
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(category); err != nil {
		helpers.SendError(w, err)
		return
	}

}

func (t *CategoryTransport) getCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		slog.Error("User id not found in context")
		helpers.SendError(w, helpers.ErrInternal)
		return
	}

	categories, err := t.service.GetCategories(ctx, userID)
	if err != nil {
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		helpers.SendError(w, err)
		return
	}
}
