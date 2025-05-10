package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
	"github.com/Corray333/keep_it/internal/errs"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/spf13/viper"
)

type logIner interface {
	storage.Transactioner
	verificationCodeGetter
	refreshTokenSetter
	userGetter
}

func (s *UserService) LogIn(ctx context.Context, user *entities.User, code string) (fullUser *entities.User, accessToken string, refreshToken string, err error) {

	query, err := s.logIner.GetCodeRequest(ctx, user.Username)
	if err != nil {
		return nil, "", "", err
	}
	if query == nil {
		return nil, "", "", helpers.NewError(http.StatusUnauthorized, errors.New("code not found"))
	}

	if query.Type != CodeRequestTypeLogIn {
		return nil, "", "", errs.ErrWrongCodeRequestType
	}

	if query.Code != code {
		return nil, "", "", errs.ErrWrongVerificationCode
	}

	fullUser, err = s.logIner.GetUser(ctx, user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to login user: ", "error", err)
	}

	if !auth.Verify(fullUser.Password, user.Password) {
		return nil, "", "", errs.ErrWrongPassword
	}

	fullUser.Password = ""

	refreshToken, err = auth.CreateToken(fullUser.ID, viper.GetDuration("auth.refresh_token_lifetime"))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create access token: ", "error", err)
	}
	creds, err := auth.ExtractCredentials(refreshToken)
	if err != nil {
		slog.Error("failed to extract credentials: ", "error", err)
		return nil, "", "", err
	}

	err = s.logIner.SetRefreshToken(ctx, fullUser.ID, refreshToken, creds.Exp)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to set refresh token: ", "error", err)
	}

	accessToken, err = auth.CreateToken(fullUser.ID, viper.GetDuration("auth.access_token_lifetime"))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create access token: ", "error", err)
	}

	return fullUser, accessToken, refreshToken, err
}
