package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
)

type userGetter interface {
	GetUser(ctx context.Context, searchUser *entities.User) (*entities.User, error)
}

func (s *UserService) GetUser(ctx context.Context, searchUser *entities.User) (*entities.User, error) {
	return s.userGetter.GetUser(ctx, searchUser)
}
