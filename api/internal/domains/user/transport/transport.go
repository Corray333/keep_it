package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type UserTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	SignUp(ctx context.Context, user entities.User, code string) (userID int64, accessToken string, refreshToken string, err error)
	LogIn(ctx context.Context, user *entities.User, code string) (fullUser *entities.User, accessToken string, refreshToken string, err error)
	RenewTokens(ctx context.Context, userID int64, oldRefreshToken string) (accessToken, refreshToken string, err error)
	FindUserByUsernameOrEmail(ctx context.Context, checkStr string) (user *entities.User, err error)
	CheckCode(ctx context.Context, testCode string, checkStr string) (correct bool, err error)
}

func New(router *chi.Mux, service service) *UserTransport {
	return &UserTransport{
		router:  router,
		service: service,
	}
}

func (t *UserTransport) RegisterRoutes() {
	t.router.Post("/api/auth/signup", t.signUp)
	t.router.Post("/api/auth/login", t.logIn)
	t.router.Post("/api/auth/renew-tokens", t.renewTokens)
	t.router.Post("/api/auth/check-code", t.checkCode)
	t.router.Post("/api/users/login-find", t.findUser)

	t.router.Group(func(r chi.Router) {
	})
}

type SignUpRequest struct {
	Password string `json:"password" example:"QWerty123"`
	Username string `json:"username" example:"corray"`
	Code     string `json:"code" example:"H78FW2"`
}

type SignUpResponse struct {
	Authorization string        `json:"authorization"`
	User          entities.User `json:"user,omitempty"`
}

func (t *UserTransport) signUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := entities.User{
		Username: req.Username,
		Password: req.Password,
	}

	userID, accessToken, refreshToken, err := t.service.SignUp(ctx, user, req.Code)
	if err != nil {
		slog.Error("failed to sign up: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = userID
	user.Password = ""

	creds, err := auth.ExtractCredentials(refreshToken)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		slog.Error("Failed to insert user: " + err.Error())
		return
	}

	cookie := http.Cookie{
		Name:     "Refresh",
		Value:    refreshToken,
		Expires:  creds.Exp,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

	respondJSON(w, SignUpResponse{
		Authorization: accessToken,
		User:          user,
	})
}

func respondJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Code     string `json:"code" example:"H78FW2"`
}

type LogInResponse struct {
	Authorization string        `json:"authorization"`
	User          entities.User `json:"user,omitempty"`
}

func (t *UserTransport) logIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &entities.User{
		Username: req.Username,
		Password: req.Password,
	}

	fullUser, accessToken, refreshToken, err := t.service.LogIn(ctx, user, req.Code)
	if err != nil {
		slog.Error("failed to log in: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = ""

	creds, err := auth.ExtractCredentials(refreshToken)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		slog.Error("Failed to insert user: " + err.Error())
		return
	}

	cookie := http.Cookie{
		Name:     "Refresh",
		Value:    refreshToken,
		Expires:  creds.Exp,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

	respondJSON(w, LogInResponse{
		Authorization: accessToken,
		User:          *fullUser,
	})
}

type RenewTokensResponse struct {
	Authorization string `json:"authorization"`
}

func (t *UserTransport) renewTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	oldRefreshToken, err := r.Cookie("Refresh")
	if err != nil {
		slog.Error("failed to get cookie: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	creds, err := auth.ExtractCredentials(oldRefreshToken.Value)
	if err != nil {
		slog.Error("failed to extract user id: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := t.service.RenewTokens(ctx, creds.ID, oldRefreshToken.Value)
	if err != nil {
		slog.Error("failed to renew tokens: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	creds, err = auth.ExtractCredentials(refreshToken)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		slog.Error("Failed to insert user: " + err.Error())
		return
	}
	cookie := http.Cookie{
		Name:     "Refresh",
		Value:    refreshToken,
		Expires:  creds.Exp,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

	respondJSON(w, LogInResponse{
		Authorization: accessToken,
	})
}

type FindUserRequest struct {
	CheckStr string `json:"checkStr"`
}

func (t *UserTransport) findUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req FindUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := t.service.FindUserByUsernameOrEmail(ctx, req.CheckStr)
	if err != nil {
		slog.Error("failed to find user: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = ""

	respondJSON(w, user)
}

type CheckCodeRequest struct {
	TestCode string `json:"testCode"`
	CheckStr string `json:"checkStr"`
}

type CheckCodeResponse struct {
	Correct bool `json:"correct"`
}

func (t *UserTransport) checkCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CheckCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	correct, err := t.service.CheckCode(ctx, req.TestCode, req.CheckStr)
	if err != nil {
		slog.Error("failed to check code: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, CheckCodeResponse{
		Correct: correct,
	})
}
