package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/spf13/viper"
)

const (
	CodeRequestTypeSignUp = iota + 1
	CodeRequestTypeLogIn
	CodeRequestTypeChangePassword
)

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	GetCodeRequest(ctx context.Context, username string) (*entities.CodeQuery, error)
	InsertUser(ctx context.Context, user entities.User) (int64, error)
	LoginUser(ctx context.Context, user *entities.User) (fullUser *entities.User, err error)
	CreateRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt int64) (err error)

	FindUserByUsernameOrEmail(ctx context.Context, checkStr string) (user *entities.User, err error)

	RenewTokens(ctx context.Context, userID int64, oldRefreshToken, newRefreshToken string, expiresAt int64) (err error)

	GetUserByTelegramID(ctx context.Context, telegramID int64) (user *entities.User, err error)
}

type UserService struct {
	repo repository
}

func New(repo repository) *UserService {
	s := &UserService{
		repo: repo,
	}
	return s
}

func (s *UserService) Run() {
}

// Errors
var (
	ErrWrongVerificationCode = helpers.NewError(http.StatusUnauthorized, errors.New("wrong verification code"))
	ErrWrongCodeRequestType  = helpers.NewError(http.StatusBadRequest, errors.New("wrong code request type"))
	ErrWrongPassword         = helpers.NewError(http.StatusUnauthorized, errors.New("wrong password"))
	ErrWrongPasswordFormat   = helpers.NewError(http.StatusBadRequest, errors.New("wrong password format"))
)

func (s *UserService) SignUp(ctx context.Context, user entities.User, code string) (userID int64, accessToken string, refreshToken string, err error) {

	query, err := s.repo.GetCodeRequest(ctx, user.Username)
	if err != nil {
		return 0, "", "", err
	}
	if query == nil {
		return 0, "", "", helpers.NewError(http.StatusUnauthorized, errors.New("code not found"))
	}

	if query.Type != CodeRequestTypeSignUp {
		return 0, "", "", ErrWrongCodeRequestType
	}

	if query.Code != code {
		return 0, "", "", ErrWrongVerificationCode
	}

	user.TelegramID = query.TelegramID

	passwordRegex := regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`)
	if !passwordRegex.MatchString(user.Password) {
		return 0, "", "", ErrWrongPasswordFormat
	}

	passHash, err := auth.Hash(user.Password)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to hash password: " + err.Error())
	}
	user.Password = passHash

	ctx, err = s.repo.Begin(ctx)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to begin transaction: " + err.Error())
	}
	defer s.repo.Rollback(ctx)

	userID, err = s.repo.InsertUser(ctx, user)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to insert user: " + err.Error())
	}

	refreshToken, err = auth.CreateToken(user.ID, viper.GetDuration("auth.refresh_token_lifetime"))
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to create access token: " + err.Error())
	}

	creds, err := auth.ExtractCredentials(refreshToken)
	if err != nil {
		slog.Error("failed to extract credentials: " + err.Error())
		return 0, "", "", err
	}

	err = s.repo.CreateRefreshToken(ctx, userID, refreshToken, creds.Exp.Unix())
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to set refresh token: " + err.Error())
	}

	err = s.repo.Commit(ctx)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to commit transaction: " + err.Error())
	}

	accessToken, err = auth.CreateToken(userID, viper.GetDuration("auth.access_token_lifetime"))
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to create access token: " + err.Error())
	}

	return userID, accessToken, refreshToken, err
}

func (s *UserService) CodeExists(ctx context.Context, username string, syn int64) (bool, error) {
	query, err := s.repo.GetCodeRequest(ctx, username)
	if err != nil {
		return false, err
	}

	if query == nil {
		return false, nil
	}

	if query.Syn != syn {
		return false, nil
	}

	return true, nil
}

func (s *UserService) LogIn(ctx context.Context, user *entities.User, code string) (fullUser *entities.User, accessToken string, refreshToken string, err error) {

	// query, err := s.repo.GetCodeRequest(ctx, user.Username)
	// if err != nil {
	// 	return nil, "", "", err
	// }
	// if query == nil {
	// 	return nil, "", "", helpers.NewError(http.StatusUnauthorized, errors.New("code not found"))
	// }

	// if query.Type != CodeRequestTypeLogIn {
	// 	return nil, "", "", ErrWrongCodeRequestType
	// }

	// if query.Code != code {
	// 	return nil, "", "", ErrWrongVerificationCode
	// }

	ctx, err = s.repo.Begin(ctx)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to begin transaction: " + err.Error())
	}
	defer s.repo.Rollback(ctx)

	fullUser, err = s.repo.LoginUser(ctx, user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to login user: " + err.Error())
	}

	if !auth.Verify(fullUser.Password, user.Password) {
		return nil, "", "", ErrWrongPassword
	}

	fullUser.Password = ""

	refreshToken, err = auth.CreateToken(fullUser.ID, viper.GetDuration("auth.refresh_token_lifetime"))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create access token: " + err.Error())
	}
	creds, err := auth.ExtractCredentials(refreshToken)
	if err != nil {
		slog.Error("failed to extract credentials: " + err.Error())
		return nil, "", "", err
	}

	err = s.repo.CreateRefreshToken(ctx, fullUser.ID, refreshToken, creds.Exp.Unix())
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to set refresh token: " + err.Error())
	}

	err = s.repo.Commit(ctx)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to commit transaction: " + err.Error())
	}

	accessToken, err = auth.CreateToken(fullUser.ID, viper.GetDuration("auth.access_token_lifetime"))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create access token: " + err.Error())
	}

	return fullUser, accessToken, refreshToken, err
}

func (s *UserService) RenewTokens(ctx context.Context, userID int64, oldRefreshToken string) (accessToken, refreshToken string, err error) {
	ctx, err = s.repo.Begin(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to begin transaction: " + err.Error())
	}
	defer s.repo.Rollback(ctx)

	refreshToken, err = auth.CreateToken(userID, viper.GetDuration("auth.refresh_token_lifetime"))
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: " + err.Error())
	}
	creds, err := auth.ExtractCredentials(refreshToken)
	if err != nil {
		slog.Error("failed to extract credentials: " + err.Error())
		return "", "", err
	}

	err = s.repo.RenewTokens(ctx, userID, oldRefreshToken, refreshToken, creds.Exp.Unix())
	if err != nil {
		return "", "", fmt.Errorf("failed to renew tokens: " + err.Error())
	}

	err = s.repo.Commit(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to commit transaction: " + err.Error())
	}

	accessToken, err = auth.CreateToken(userID, viper.GetDuration("auth.access_token_lifetime"))
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: " + err.Error())
	}

	return accessToken, refreshToken, err
}

func (s *UserService) FindUserByUsernameOrEmail(ctx context.Context, checkStr string) (user *entities.User, err error) {
	return s.repo.FindUserByUsernameOrEmail(ctx, checkStr)
}

func (s *UserService) GetUserByTelegramID(ctx context.Context, telegramID int64) (user *entities.User, err error) {
	return s.repo.GetUserByTelegramID(ctx, telegramID)
}
