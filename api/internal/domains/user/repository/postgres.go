package repository

import (
	"context"
	"log/slog"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
)

func (s *UserRepository) InsertUser(ctx context.Context, user entities.User) (int64, error) {

	tx, isNew, err := s.GetTx(ctx)
	if err != nil {
		slog.Error("failed to get transaction: " + err.Error())
		return -1, err
	}
	if isNew {
		defer tx.Rollback()
	}

	rows := tx.QueryRow(`
		INSERT INTO users (username, email, tg_username, password, avatar, ref_code) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id;
	`, user.Username, user.Email, user.TelegramUsername, user.Password, "/images/avatars/default_avatar.png", user.RefCode)

	if err := rows.Scan(&user.ID); err != nil {
		return -1, err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return -1, err
		}
	}

	return user.ID, nil
}

func (s *UserRepository) LoginUser(ctx context.Context, user entities.User) (userID int64, correctPassword string, err error) {

	tx, isNew, err := s.GetTx(ctx)
	if err != nil {
		return -1, "", err
	}
	if isNew {
		defer tx.Rollback()
	}

	rows := s.DB.QueryRow(`
		SELECT user_id, password FROM users WHERE email = $1;
	`, user.Email)

	if err := rows.Scan(&userID, &correctPassword); err != nil {
		return -1, "", err
	}

	return user.ID, correctPassword, nil
}

func (s *UserRepository) CreateRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt int64) (err error) {
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

func (s *UserRepository) RenewTokens(ctx context.Context, userID int64, oldRefreshToken, newRefreshToken string, expiresAt int64) (err error) {
	tx, isNew, err := s.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	_, err = tx.Exec(`UPDATE user_token SET token = $1, expires_at = $2 WHERE user_id = $3 AND token = $4;`, newRefreshToken, expiresAt, userID, oldRefreshToken)
	if err != nil {
		slog.Error("error updating refresh token: " + err.Error())
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
