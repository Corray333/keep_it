package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
)

func (s *UserRepository) FindUserByUsernameOrEmail(ctx context.Context, checkStr string) (user *entities.User, err error) {
	user = &entities.User{}

	if err = s.DB.Get(user, `
		SELECT * FROM users WHERE username = $1 OR email = $1;
	`, checkStr); err != nil {
		return nil, err
	}

	return &entities.User{
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

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
		INSERT INTO users (username, email, tg_id, password, avatar, ref_code) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id;
	`, user.Username, user.Email, user.TelegramID, user.Password, "/images/avatars/default_avatar.png", user.RefCode)

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

func (s *UserRepository) LoginUser(ctx context.Context, user *entities.User) (fullUser *entities.User, err error) {

	tx, isNew, err := s.GetTx(ctx)
	if err != nil {
		return nil, err
	}
	if isNew {
		defer tx.Rollback()
	}

	fullUser = &entities.User{}

	if err = s.DB.Get(fullUser, `
		SELECT * FROM users WHERE username = $1;
	`, user.Username); err != nil {
		return nil, err
	}

	return fullUser, nil
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

func (s *UserRepository) GetUserByTelegramID(ctx context.Context, telegramID int64) (user *entities.User, err error) {
	user = &entities.User{}

	if err = s.DB.Get(user, `
		SELECT * FROM users WHERE tg_id = $1;
	`, telegramID); err != nil {
		return nil, err
	}

	return user, nil
}
