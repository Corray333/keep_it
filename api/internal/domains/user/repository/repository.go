package repository

import (
	"github.com/Corray333/keep_it/internal/storage"
)

type UserRepository struct {
	*storage.Storage
}

func New(store *storage.Storage) *UserRepository {
	return &UserRepository{
		Storage: store,
	}
}
