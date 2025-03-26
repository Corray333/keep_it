package repository

import (
	"context"

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

func (s *UserRepository) GetUserByTelegramID(ctx context.Context, telegramID int64) (user *entities.User, err error) {
	user = &entities.User{}

	if err = s.DB.Get(user, `
		SELECT * FROM users WHERE tg_id = $1;
	`, telegramID); err != nil {
		return nil, err
	}

	return user, nil
}
