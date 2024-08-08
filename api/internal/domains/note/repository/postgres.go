package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Corray333/keep_it/internal/domains/note/types"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
)

type NoteRepository struct {
	db *sqlx.DB
}

var (
	ErrorNoteDoesNotExist        error = errors.New("note does not exist")
	ErrorTagsFieldHasInvalidData error = errors.New("tags field has invalid data")
)

// New creates a new storage and tables
func NewStorage(db *sqlx.DB, redis *redis.Client) *NoteRepository {
	store := &NoteRepository{
		db: db,
	}

	return store
}

func (s *NoteRepository) CreateNote(note *types.Note) (string, error) {

	// TODO: forbid if no access
	tx, err := s.db.Beginx()
	defer tx.Rollback()
	if err != nil {
		return "", err
	}

	note_id := ""

	if err := tx.QueryRow("INSERT INTO notes (creator, title, source, original, font, created_at, type, category_owner, category_id, content, icon) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING note_id", note.Creator, note.Title, note.Source, note.OriginalRaw, note.Font, note.CreatedAt, note.Type, note.CategoryOwner, note.CategoryId, note.ContentRaw, note.IconRaw).Scan(&note_id); err != nil {
		return "", err
	}

	if _, err := tx.Exec("INSERT INTO user_note_access VALUES($1, $2)", note.Creator, note_id); err != nil {
		return "", err
	}

	return note_id, tx.Commit()
}

func (s *NoteRepository) GetNote(note_id string) (*types.Note, error) {
	// TODO: forbid if no access
	note := &types.Note{}
	rows := s.db.QueryRowx("SELECT * FROM notes WHERE note_id = $1", note_id)
	if err := rows.StructScan(note); err != nil {
		return nil, err
	}
	if err := s.db.Select(&note.Tags, "SELECT tag_id, owner, text, color FROM note_tag NATURAL JOIN tags WHERE note_id = $1", note_id); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(note.ContentRaw), &note.Content); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(*note.OriginalRaw, &note.Original); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(note.IconRaw), &note.Icon); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteRepository) GetNotes(user_id int, offset int, filter map[string]interface{}) ([]*types.Note, bool, error) {
	sqfilter := sq.Eq(filter)
	sqfilter["user_id"] = user_id

	// TODO: add pagination

	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("notes.note_id, creator, title, source, original, font, created_at, copied_at, type, checked, content, icon, cover, category_owner, category_id").
		From("user_note_access").
		Join("notes on user_note_access.note_id = notes.note_id").
		Where(sqfilter).
		Offset(uint64(offset)).
		Limit(51).
		ToSql()
	if err != nil {
		return nil, false, err
	}

	notes := []*types.Note{}

	if err := s.db.Select(&notes, sql, args...); err != nil {
		return nil, false, err
	}

	for _, note := range notes {
		if err := s.db.Select(&note.Tags, "SELECT tag_id, owner, text, color FROM note_tag NATURAL JOIN tags WHERE note_id = $1", note.ID); err != nil {
			return nil, false, err
		}

		if err := json.Unmarshal([]byte(note.ContentRaw), &note.Content); err != nil {
			return nil, false, err
		}

		if err := json.Unmarshal(*note.OriginalRaw, &note.Original); err != nil {
			return nil, false, err
		}

		if err := json.Unmarshal([]byte(note.IconRaw), &note.Icon); err != nil {
			return nil, false, err
		}
	}
	if len(notes) == 51 {
		return notes[:50], true, nil
	}

	return notes, false, nil
}

func (s *NoteRepository) CheckNoteAccess(note_id string, user_id int) (bool, error) {
	// TODO: forbid if no access
	exists := false
	err := s.db.QueryRowx("SELECT EXISTS(SELECT 1 FROM user_note_access WHERE note_id = $1 AND user_id = $2)", note_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil

}

func (s *NoteRepository) CreateTag(tag *types.Tag) (*types.Tag, error) {
	err := s.db.QueryRowx("INSERT INTO tags (text, color, owner) VALUES ($1, $2, $3) RETURNING tag_id", tag.Text, tag.Color, tag.Owner).Scan(&tag.ID)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *NoteRepository) UpdateNote(note_id string, data map[string]interface{}) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tagsRaw, ok := data["tags"]
	if ok {
		if _, err := tx.Exec("DELETE FROM note_tag WHERE note_id = $1", note_id); err != nil {
			return err
		}
		query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("note_tag")
		fmt.Printf("%T\n", tagsRaw)
		tags, ok := tagsRaw.([]interface{})
		if !ok {
			fmt.Println("not a slice")
			return ErrorTagsFieldHasInvalidData
		}
		for _, tagRaw := range tags {
			tag, ok := tagRaw.(map[string]interface{})
			if !ok {
				fmt.Println("not a map")
				return ErrorTagsFieldHasInvalidData
			}
			tagID, ok := tag["id"]
			if !ok {
				fmt.Println("no id")
				return ErrorTagsFieldHasInvalidData
			}
			owner, ok := tag["owner"]
			if !ok {
				fmt.Println("no owner")
				return ErrorTagsFieldHasInvalidData
			}
			query = query.Values(note_id, tagID, owner)
		}
		sql, args, err := query.ToSql()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(sql, args...); err != nil {
			return err
		}
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update("notes")
	if len(maps.Keys(data)) == 1 && ok {
		return tx.Commit()
	}
	for key, val := range data {
		if key == "tags" {
			continue
		}
		query = query.Set(key, val)
	}
	query = query.Where(sq.Eq{"note_id": note_id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	res, err := tx.Exec(sql, args...)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return ErrorNoteDoesNotExist
	}

	return tx.Commit()
}

func (s *NoteRepository) DeleteNote(note_id string, uid int) error {

	res, err := s.db.Exec("DELETE FROM notes WHERE note_id = $1 AND creator = $2", note_id, uid)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
