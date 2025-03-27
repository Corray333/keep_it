package repository

import (
	"github.com/Corray333/keep_it/internal/storage"
)

type NoteRepository struct {
	*storage.Storage
}

func New(store *storage.Storage) *NoteRepository {
	return &NoteRepository{
		Storage: store,
	}
}
