package repository

import (
	"context"
	"log/slog"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
)

func (s *UserRepository) NewUser(ctx context.Context, user entities.User) (int64, error) {

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
