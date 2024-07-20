package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Corray333/keep_it/internal/domains/user/types"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/rand"
)

const MaxFileSize = 5 << 20

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

type SignUpResponse struct {
	Authorization string     `json:"authorization"`
	User          types.User `json:"user,omitempty"`
}

type SignUpRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Code     string `json:"code"`
}

// SignUp handles user signup
// @Summary Sign up a new user
// @Description Sign up a new user with username, password, and verification code
// @Tags users
// @Accept json
// @Produce json
// @Param signupRequest body SignUpRequest true "Sign up request"
// @Success 200 {object} SignUpResponse
// @Failure 400 {string} string "Bad request"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /api/users/signup [post]
func SignUp(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := SignUpRequest{}
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
		user := types.User{
			Avatar:           "/images/avatars/default_avatar.png",
			Username:         request.Username,
			Password:         request.Password,
			TelegramUsername: query.TG,
		}

		id, refresh, err := store.InsertUser(user, r.UserAgent())
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			slog.Error("Failed to insert user: " + err.Error())
			return
		}
		user.ID = id

		creds, err := auth.ExtractCredentials(refresh)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			slog.Error("Failed to insert user: " + err.Error())
			return
		}

		cookie := http.Cookie{
			Name:     "Refresh",
			Value:    refresh,
			Expires:  creds.Exp,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}

		http.SetCookie(w, &cookie)

		token, err := auth.CreateToken(user.ID, auth.AccessTokenLifeTime)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			slog.Error("Failed to create token: " + err.Error())
			return
		}
		user.Password = ""
		if err := json.NewEncoder(w).Encode(SignUpResponse{
			Authorization: token,
			User:          user,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to send response: " + err.Error())
			return
		}
	}
}

type LogInResponse struct {
	Authorization string     `json:"authorization"`
	User          types.User `json:"user,omitempty"`
}
type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// LogIn handles user login
// @Summary Log in a user
// @Description Log in a user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login request"
// @Success 200 {object} LogInResponse
// @Failure 400 {string} string "Bad request"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /api/users/login [post]
func LogIn(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: check and remove expired tokens

		request := LoginRequest{}
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
		user := types.User{
			Username: request.Username,
			Password: request.Password,
		}
		id, refresh, err := store.LoginUser(user, r.UserAgent())
		if err != nil {
			http.Error(w, "Wrong password or email", http.StatusForbidden)
			slog.Error("Failed to login user: " + err.Error())
			return
		}
		user.ID = id

		fmt.Println()
		fmt.Println("Login refresh: ", refresh)
		fmt.Println()

		creds, err := auth.ExtractCredentials(refresh)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			slog.Error("Failed to insert user: " + err.Error())
			return
		}

		cookie := http.Cookie{
			Name:     "Refresh",
			Value:    refresh,
			Expires:  creds.Exp,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}

		http.SetCookie(w, &cookie)

		token, err := auth.CreateToken(user.ID, auth.AccessTokenLifeTime)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			slog.Error("Failed to create token: " + err.Error())
			return
		}
		user.Password = ""
		if err := json.NewEncoder(w).Encode(LogInResponse{
			Authorization: token,
			User:          user,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to send response: " + err.Error())
			return
		}
	}
}

type RefreshAccessTokenResponse struct {
	Authorization string `json:"authorization"`
}

// RefreshAccessToken handles refreshing of access token
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token in cookies
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} RefreshAccessTokenResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/users/refresh [get]
func RefreshAccessToken(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshCookie, err := r.Cookie("Refresh")
		if err != nil {
			http.Error(w, "Failed to get refresh cookie", http.StatusUnauthorized)
			slog.Error("Failed to get refresh cookie: " + err.Error())
			return
		}
		if refreshCookie.Value == "" {
			http.Error(w, "Failed to get refresh cookie", http.StatusUnauthorized)
			slog.Error("Failed to get refresh cookie")
			return
		}

		fmt.Println()
		fmt.Println("Cookie: ", refreshCookie.Value)
		fmt.Println()

		creds, err := auth.ExtractCredentials(refreshCookie.Value)
		if err != nil {
			http.Error(w, "Failed to extract credentials", http.StatusBadRequest)
			slog.Error("Failed to extract credentials: " + err.Error())
			return
		}
		access, refresh, err := store.RefreshToken(creds.ID, r.UserAgent(), refreshCookie.Value)
		if err != nil {
			http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
			slog.Error("Failed to refresh token: " + err.Error())
			return
		}

		creds, err = auth.ExtractCredentials(refresh)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			slog.Error("Failed to insert user: " + err.Error())
			return
		}

		cookie := http.Cookie{
			Name:     "Refresh",
			Value:    refresh,
			Expires:  creds.Exp,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}

		http.SetCookie(w, &cookie)

		if err := json.NewEncoder(w).Encode(RefreshAccessTokenResponse{
			Authorization: access,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response: " + err.Error())
			return
		}
	}
}

// GetUser handles fetching user information
// @Summary Get user information
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} types.User
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/users/{id} [get]
func GetUser(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "id")
		if userId == "0" {
			creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Failed to extract credentials", http.StatusBadRequest)
				slog.Error("Failed to extract credentials: " + err.Error())
				return
			}
			userId = strconv.Itoa(int(creds.ID))
		}

		user, err := store.SelectUser(userId)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			slog.Error("Failed to get user: " + err.Error())
			return
		}

		// TODO: create this struct
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

// UpdateUser handles updating user information
// @Summary Update user information
// @Description Update user's avatar and username
// @Tags users
// @Accept multipart/form-data
// @Param avatar formData file false "User's avatar"
// @Param username formData string true "User's username"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
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

// CheckUsername handles checking if a username is available
// @Summary Check username availability
// @Description Check if a username is available
// @Tags users
// @Produce json
// @Param username query string true "Username to check"
// @Success 200 {object} CheckUsernameResponse
// @Failure 500 {string} string "Internal server error"
// @Router /api/users/check-username [get]
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

// CheckCode handles checking if a verification code is valid
// @Summary Check verification code
// @Description Check if a verification code is valid for the given username
// @Tags users
// @Accept json
// @Produce json
// @Param checkCodeRequest body CheckCodeRequest true "Check code request"
// @Success 200 {object} CheckCodeResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/users/check-code [post]
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
			w.WriteHeader(http.StatusUnauthorized)
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

func UploadImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: add image compression
		if err := r.ParseMultipartForm(MaxFileSize); err != nil {
			slog.Error("error parsing multipart form: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			slog.Error("error getting file: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		randomStr := ""
		for i := 0; i < 10; i++ {
			randomStr += strconv.Itoa(rand.Intn(10))
		}
		fileName := randomStr + filepath.Ext(header.Filename)
		newFile, err := os.Create("../files/images/" + fileName)
		if err != nil {
			slog.Error("error creating file: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer newFile.Close()
		if _, err := io.Copy(newFile, file); err != nil {
			slog.Error("error copying file: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(struct {
			URL string `json:"url"`
		}{
			URL: "/images/" + fileName,
		}); err != nil {
			slog.Error("error encoding or sending file name: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

	}
}
