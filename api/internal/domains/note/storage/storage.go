package storage

import (
	"fmt"

	"github.com/Corray333/keep_it/internal/domains/note/types"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

type NoteStorage struct {
	db *sqlx.DB
}

// New creates a new storage and tables
func NewStorage(db *sqlx.DB, redis *redis.Client) *NoteStorage {
	return &NoteStorage{
		db: db,
	}
}

func (s *NoteStorage) CreateNote(note *types.Note) (string, error) {
	// TODO: forbid if no access
	tx, err := s.db.Beginx()
	defer tx.Rollback()
	if err != nil {
		return "", err
	}

	note_id := ""

	fmt.Println()
	fmt.Println(string(note.Content))
	fmt.Println()

	if err := tx.QueryRow("INSERT INTO notes (creator, title, source, original, font, created_at, type, category_owner, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING note_id", note.Creator, note.Title, note.Source, note.Original, note.Font, note.CreatedAt, note.Type, note.CategoryOwner, note.CategoryId).Scan(&note_id); err != nil {
		return "", err
	}

	return note_id, nil

}

func (s *NoteStorage) GetNote(note_id string) (*types.Note, error) {
	// TODO: forbid if no access
	note := &types.Note{}
	rows := s.db.QueryRowx("SELECT * FROM notes WHERE note_id = $1", note_id)
	if err := rows.StructScan(note); err != nil {
		return nil, err
	}

	return note, nil
}
