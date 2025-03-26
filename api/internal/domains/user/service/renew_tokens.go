package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Corray333/keep_it/internal/storage"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/spf13/viper"
)

type tokensRenewer interface {
	storage.Transactioner

	RenewTokens(ctx context.Context, userID int64, oldRefreshToken, newRefreshToken string, expiresAt int64) (err error)
}

func (s *UserService) RenewTokens(ctx context.Context, userID int64, oldRefreshToken string) (accessToken, refreshToken string, err error) {
	ctx, err = s.tokensRenewer.Begin(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to begin transaction: " + err.Error())
	}
	defer s.tokensRenewer.Rollback(ctx)

	refreshToken, err = auth.CreateToken(userID, viper.GetDuration("auth.refresh_token_lifetime"))
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: " + err.Error())
	}
	creds, err := auth.ExtractCredentials(refreshToken)
	if err != nil {
		slog.Error("failed to extract credentials: " + err.Error())
		return "", "", err
	}

	err = s.tokensRenewer.RenewTokens(ctx, userID, oldRefreshToken, refreshToken, creds.Exp.Unix())
	if err != nil {
		return "", "", fmt.Errorf("failed to renew tokens: " + err.Error())
	}

	err = s.tokensRenewer.Commit(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to commit transaction: " + err.Error())
	}

	accessToken, err = auth.CreateToken(userID, viper.GetDuration("auth.access_token_lifetime"))
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: " + err.Error())
	}

	return accessToken, refreshToken, err
}
