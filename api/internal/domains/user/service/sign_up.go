package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
	"github.com/Corray333/keep_it/internal/errs"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/spf13/viper"
)

type signUper interface {
	storage.Transactioner
	verificationCodeGetter
	refreshTokenSetter

	NewUser(ctx context.Context, user entities.User) (int64, error)
}

type verificationCodeGetter interface {
	GetCodeRequest(ctx context.Context, username string) (*entities.CodeQuery, error)
}

type refreshTokenSetter interface {
	SetRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt time.Time) (err error)
}

func (s *UserService) SignUp(ctx context.Context, user entities.User, code string) (userID int64, accessToken string, refreshToken string, err error) {

	query, err := s.signUper.GetCodeRequest(ctx, user.Username)
	if err != nil {
		return 0, "", "", err
	}
	if query == nil {
		return 0, "", "", helpers.NewError(http.StatusUnauthorized, errors.New("code not found"))
	}

	if query.Type != CodeRequestTypeSignUp {
		return 0, "", "", errs.ErrWrongCodeRequestType
	}

	if query.Code != code {
		return 0, "", "", errs.ErrWrongVerificationCode
	}

	user.TelegramID = query.TelegramID

	passwordRegex := regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`)
	if !passwordRegex.MatchString(user.Password) {
		return 0, "", "", errs.ErrWrongPasswordFormat
	}

	passHash, err := auth.Hash(user.Password)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to hash password: ", "error", err)
	}
	user.Password = passHash

	ctx, err = s.signUper.Begin(ctx)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to begin transaction: ", "error", err)
	}
	defer s.signUper.Rollback(ctx)

	userID, err = s.signUper.NewUser(ctx, user)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to insert user: ", "error", err)
	}

	refreshToken, err = auth.CreateToken(userID, viper.GetDuration("auth.refresh_token_lifetime"))
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to create access token: ", "error", err)
	}

	creds, err := auth.ExtractCredentials(refreshToken)
	if err != nil {
		slog.Error("failed to extract credentials: ", "error", err)
		return 0, "", "", err
	}

	err = s.signUper.SetRefreshToken(ctx, userID, refreshToken, creds.Exp)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to set refresh token: ", "error", err)
	}

	err = s.signUper.Commit(ctx)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to commit transaction: ", "error", err)
	}

	accessToken, err = auth.CreateToken(userID, viper.GetDuration("auth.access_token_lifetime"))
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to create access token: ", "error", err)
	}

	return userID, accessToken, refreshToken, err
}
