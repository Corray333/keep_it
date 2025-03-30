package repository

import (
	"context"
	"log/slog"
)

func (r *UserRepository) SetVKID(ctx context.Context, userID int64, vkID int64) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	_, err = tx.Exec(`UPDATE users SET vk_id = NULL WHERE vk_id = $1;`, vkID)
	if err != nil {
		slog.Error("error setting refresh token: " + err.Error())
		return err
	}

	_, err = tx.Exec(`UPDATE users SET vk_id = $1 WHERE user_id = $2;`, vkID, userID)
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
