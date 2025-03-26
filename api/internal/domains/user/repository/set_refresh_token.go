package repository

import (
	"context"
	"log/slog"
)

func (s *UserRepository) SetRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt int64) (err error) {
	tx, isNew, err := s.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	_, err = tx.Exec(`INSERT INTO user_token (user_id, token, expires_at) VALUES ($1, $2, $3);`, userID, refreshToken, expiresAt)
	if err != nil {
		slog.Error("error setting refresh token: " + err.Error())
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
