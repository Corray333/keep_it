package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenLifeTime  = time.Minute * 15
	RefreshTokenLifeTime = time.Hour * 24 * 365
)

var secretKey []byte

// init initializes the secret key from the environment variable
func init() {
	secretKey = []byte(os.Getenv("SECRET_KEY"))
}

func NewAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		slog.Info("auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			creds, err := VerifyToken(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				slog.Error("Unauthorized: " + err.Error())
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), "creds", creds))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Hash hashes the password using bcrypt package
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Verify checks if the hashed password is equal to the password using bcrypt package
func Verify(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

// CreateToken creates a new JWT token by the email
func CreateToken(id int64, lifeTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(lifeTime).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken checks if the JWT is valid
func VerifyToken(tokenString string) (Credentials, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return Credentials{}, err
	}

	if !token.Valid {
		return Credentials{}, fmt.Errorf("invalid token")
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return Credentials{}, err
	}
	creds := Credentials{
		ID:  int64(claims["id"].(float64)),
		Exp: exp.Time,
	}

	return creds, nil
}

type Credentials struct {
	Email string `json:"email"`
	ID    int64  `json:"id,omitempty" db:"user_id"`
	Exp   time.Time
}

func ExtractCredentials(tokenString string) (*Credentials, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	credentials := Credentials{
		ID:  int64(claims["id"].(float64)),
		Exp: exp.Time,
	}
	return &credentials, nil
}
