package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/Corray333/keep_it/internal/domains/user/types"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type Storage interface {
	InsertUser(user types.User, agent string) (int, string, error)
	LoginUser(user types.User, agent string) (int, string, error)
	RefreshToken(id int, agent string, refresh string) (string, string, error)
	SelectUser(id string) (types.User, error)
	UpdateUser(user types.User) error
	CheckUsername(username string) (bool, error)
	GetCodeRequest(username string) (*types.CodeQuery, error)
}

const (
	CodeRequestTypeSignUp = iota + 1
	CodeRequestTypeLogIn
	CodeRequestTypeChangePassword
)

func RequestCodeByEmail(store *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

type LogInResponse struct {
	Authorization string     `json:"authorization"`
	Refresh       string     `json:"refresh"`
	User          types.User `json:"user,omitempty"`
}

type SignUpRequest struct {
	types.User
	Code string `json:"code"`
}

// SignUp registers a new user and returns the user ID, refresh token, and access token.
// @Summary Sign up a new user
// @Description Registers a new user and returns the user ID, refresh token, and access token.
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body types.User true "User details"
// @Success 200 {object} LogInResponse
// @Router /api/users//signup [post]
func SignUp(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := SignUpRequest{}
		request.Avatar = "/images/avatars/default_avatar.png"
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			slog.Error("Failed to read request body: " + err.Error())
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
			slog.Error("Failed to unmarshal request body: " + err.Error())
			return
		}
		query, err := store.GetCodeRequest(request.Username)
		if err != nil {
			http.Error(w, "Failed to find verification code: ", http.StatusInternalServerError)
			slog.Error("Failed to find verification code: " + err.Error())
			return
		}
		if query.TypeID != CodeRequestTypeSignUp {
			http.Error(w, "Wrong request type", http.StatusBadRequest)
			slog.Error("Wrong type of code request: has to be sign up (1)")
			return
		}
		if query.Code != request.Code {
			http.Error(w, "Wrong verification code", http.StatusForbidden)
			return
		}
		user := request.User
		user.TelegramUsername = query.TG

		id, refresh, err := store.InsertUser(user, r.UserAgent())
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			slog.Error("Failed to insert user: " + err.Error())
			return
		}
		user.ID = id

		token, err := auth.CreateToken(user.ID, auth.AccessTokenLifeTime)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			slog.Error("Failed to create token: " + err.Error())
			return
		}
		user.Password = ""
		if err := json.NewEncoder(w).Encode(LogInResponse{Authorization: token,
			Refresh: refresh,
			User:    user,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to send response: " + err.Error())
			return
		}
	}
}

// LogIn logs in a user and returns the user ID, refresh token, and access token.
// @Summary Log in a user
// @Description Logs in a user and returns the user ID, refresh token, and access token.
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body types.User true "User details"
// @Success 200 {object} LogInResponse
// @Router /api/users/login [post]
func LogIn(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := types.User{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			slog.Error("Failed to read request body: " + err.Error())
			return
		}
		if err := json.Unmarshal(body, &user); err != nil {
			http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
			slog.Error("Failed to unmarshal request body: " + err.Error())
			return
		}
		id, refresh, err := store.LoginUser(user, r.UserAgent())
		if err != nil {
			http.Error(w, "Wrong password or email", http.StatusForbidden)
			slog.Error("Failed to login user: " + err.Error())
			return
		}
		user.ID = id

		token, err := auth.CreateToken(user.ID, auth.AccessTokenLifeTime)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			slog.Error("Failed to create token: " + err.Error())
			return
		}
		user.Password = ""
		if err := json.NewEncoder(w).Encode(LogInResponse{
			Authorization: token,
			Refresh:       refresh,
			User:          user,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to send response: " + err.Error())
			return
		}
	}
}

// RefreshAccessToken refreshes the access token using the refresh token.
// @Summary Refresh access token
// @Description Refreshes the access token using the refresh token.
// @Tags user
// @Accept  json
// @Produce  json
// @Param refresh query string true "Refresh token"
// @Success 200 {object} LogInResponse
// @Router /api/users/refresh_token [get]
func RefreshAccessToken(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refresh := r.Header.Get("Refresh")
		creds, err := auth.ExtractCredentials(r.Header.Get("Refresh"))
		if err != nil {
			http.Error(w, "Failed to extract credentials", http.StatusBadRequest)
			slog.Error("Failed to extract credentials: " + err.Error())
			return
		}
		access, refresh, err := store.RefreshToken(creds.ID, r.UserAgent(), refresh)
		if err != nil {
			http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
			slog.Error("Failed to refresh token: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(LogInResponse{
			Authorization: access,
			Refresh:       refresh,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response: " + err.Error())
			return
		}
	}
}

// GetUser retrieves a user by their ID.
// @Summary Get a user by ID
// @Description Retrieves a user by their ID.
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} types.User
// @Router /api/users/{id} [get]
func GetUser(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "id")
		user, err := store.SelectUser(userId)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			slog.Error("Failed to get user: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(struct {
			User types.User `json:"user"`
		}{
			User: user,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response: " + err.Error())
			return
		}
	}
}

// UpdateUser updates a user's details.
// @Summary Update a user
// @Description Updates a user's details.
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body types.User true "Updated user details"
// @Success 200 {object} types.User
// @Router /api/users/update [put]
func UpdateUser(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Failed to extract credentials", http.StatusBadRequest)
			slog.Error("Failed to extract credentials: " + err.Error())
			return
		}
		user, err := store.SelectUser(strconv.Itoa(int(creds.ID)))
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			slog.Error("Failed to get user: " + err.Error())
			return
		}
		file, _, err := r.FormFile("avatar")
		if err != nil && err.Error() != "http: no such file" {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			slog.Error("Failed to read file: " + err.Error())
			return
		}
		if file != nil {
			newFile, err := os.Create("../files/images/avatars/avatar" + strconv.Itoa(int(user.ID)) + ".png")
			if err != nil {
				http.Error(w, "Failed to create file", http.StatusInternalServerError)
				slog.Error("Failed to create file: " + err.Error())
				return
			}
			data, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Failed to read file", http.StatusInternalServerError)
				slog.Error("Failed to read file: " + err.Error())
				return
			}
			if _, err := newFile.Write(data); err != nil {
				http.Error(w, "Failed to write file", http.StatusInternalServerError)
				slog.Error("Failed to write file: " + err.Error())
				return
			}
			user.Avatar = "images/avatars/avatar" + strconv.Itoa(int(user.ID)) + ".png"
		}
		user.Username = r.FormValue("username")
		if err := store.UpdateUser(user); err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			slog.Error("Failed to update user: " + err.Error())
			return
		}
	}
}

type CheckUsernameResponse struct {
	Found bool `json:"found"`
}

func CheckUsername(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		res, err := store.CheckUsername(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("error during checking username in db: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(CheckUsernameResponse{
			Found: res,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("error during marshalling check username response: " + err.Error())
			return
		}
	}
}

type CheckCodeRequest struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}
type CheckCodeResponse struct {
	Valid bool `json:"valid"`
}

func CheckCode(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &CheckCodeRequest{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("error during unmarshalling check code request: " + err.Error())
			return
		}

		codeQuery, err := store.GetCodeRequest(req.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("error during checking code in db: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(CheckCodeResponse{
			Valid: req.Code == codeQuery.Code,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("error during marshalling check code response: " + err.Error())
			return
		}
	}
}
