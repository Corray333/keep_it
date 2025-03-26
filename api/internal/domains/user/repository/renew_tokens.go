package repository

import (
	"context"
	"errors"
	"log/slog"
)

func (s *UserRepository) RenewTokens(ctx context.Context, userID int64, oldRefreshToken, newRefreshToken string, expiresAt int64) (err error) {
	tx, isNew, err := s.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	res, err := tx.Exec(`UPDATE user_token SET token = $1, expires_at = $2 WHERE user_id = $3 AND token = $4;`, newRefreshToken, expiresAt, userID, oldRefreshToken)
	if err != nil {
		slog.Error("error updating refresh token: " + err.Error())
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		slog.Error("error while getting rows affected updating tokens: " + err.Error())
		return err
	}

	if affected != 1 {
		slog.Error("number of rows affected by updating tokens is not 1")
		return errors.New("number of rows affected by updating tokens is not 1")
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
