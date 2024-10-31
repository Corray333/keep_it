package service

import "context"

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
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
