package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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
		slog.Error("Error building query", "error", err)
		return nil, err
	}

	notesInternal := []noteDb{}
	if err := r.DB.Select(&notesInternal, queryStr, args...); err != nil && errors.Is(err, sql.ErrNoRows) {
		slog.Error("Error getting notes", "error", err)
		return nil, err
	}

	noteIDs := []uuid.UUID{}

	for _, note := range notesInternal {
		notes = append(notes, *note.ToNote())
		noteIDs = append(noteIDs, note.NoteID)
	}

	tags := []tagDB{}
	if len(noteIDs) > 0 {
		placeholders := make([]string, len(noteIDs))
		args := make([]interface{}, len(noteIDs))
		for i, id := range noteIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args[i] = id
		}
		query := fmt.Sprintf("SELECT note_id, tag_text, tag_color FROM note_tag NATURAL JOIN tags WHERE note_id IN (%s)", strings.Join(placeholders, ","))
		if err := r.DB.Select(&tags, query, args...); err != nil {
			slog.Error("Error getting tags", "error", err)
			return nil, err
		}
	} else {
		tags = []tagDB{}
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
