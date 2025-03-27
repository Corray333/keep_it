package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *NoteRepository) GetNotes(ctx context.Context, userID int64, offset int, filters entities.NoteFilter) ([]entities.Note, error) {
	notes := []entities.Note{}

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").From("notes").Where(squirrel.Eq{"creator_id": userID}).OrderBy("created_at DESC").Offset(uint64(offset)).Limit(20)

	if filters.Category != "" {
		sq = sq.Where(squirrel.Eq{"category_id": filters.Category})
	}

	if len(filters.Tags) > 0 {
		sq = sq.Join("note_tag ON notes.note_id = note_tag.note_id").Where(squirrel.Eq{"note_tag.tag_text": filters.Tags})
	}

	if filters.Text != "" {
		searchText := fmt.Sprintf("%%%s%%", filters.Text)
		sq = sq.Where(squirrel.Or{
			squirrel.Like{"title": searchText},
			squirrel.Like{"content": searchText},
		})
	}

	if filters.Source != "" {
		sq = sq.Where(squirrel.Eq{"source": filters.Source})
	}

	queryStr, args, err := sq.ToSql()
	if err != nil {
		return nil, err
	}

	notesInternal := []noteDb{}
	if err := r.DB.Select(&notesInternal, queryStr, args...); err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	noteIDs := []uuid.UUID{}

	for _, note := range notesInternal {
		notes = append(notes, *note.ToNote())
		noteIDs = append(noteIDs, note.NoteID)
	}

	tags := []tagDB{}
	if err := r.DB.Select(&tags, "SELECT * FROM note_tag NATURAL JOIN tags WHERE note_id IN ($1)", noteIDs); err != nil {
		return nil, err
	}

	tagsMap := map[uuid.UUID][]entities.Tag{}
	for _, tag := range tags {
		tagsMap[tag.NoteID] = append(tagsMap[tag.NoteID], *tag.ToTag())
	}

	for i := range notes {
		notes[i].Tags = tagsMap[notes[i].ID]
		if notes[i].Tags == nil {
			notes[i].Tags = []entities.Tag{}
		}
	}

	return notes, nil
}
