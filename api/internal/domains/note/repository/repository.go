package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/Masterminds/squirrel"
)

type NoteRepository struct {
	*storage.Storage
}

func New(store *storage.Storage) *NoteRepository {
	return &NoteRepository{
		Storage: store,
	}
}

func (r *NoteRepository) GetNote(ctx context.Context, noteID string) (*entities.Note, error) {
	note := &entities.Note{}
	if err := r.DB.Get(note, "SELECT * FROM notes WHERE note_id = $1", noteID); err != nil {
		return nil, err
	}

	// TODO: get category and tags

	return note, nil
}

func (r *NoteRepository) CreateNote(ctx context.Context, note *entities.Note) (noteID string, err error) {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return "", err
	}
	if isNew {
		defer tx.Rollback()
	}

	fmt.Println("Created at: ", note.CreatedAt)

	if err := tx.QueryRow("INSERT INTO notes (creator_id, title, source, original, created_at, type, category_id, content, icon, cover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING note_id", note.CreatorID, note.Title, note.Source, note.Original, note.CreatedAt, note.Type, note.CategoryId, note.Content, note.Icon, note.Cover).Scan(&noteID); err != nil {
		return "", err
	}

	// TODO: insert tags and access

	if isNew {
		if err := tx.Commit(); err != nil {
			return "", err
		}
	}

	return noteID, nil
}

type GetNotesInternal struct {
	NoteID     string          `db:"note_id"`
	CreatorID  int64           `db:"creator_id"`
	Title      string          `db:"title"`
	Source     string          `db:"source"`
	Original   string          `db:"original"`
	Icon       json.RawMessage `db:"icon"`
	CreatedAt  int64           `db:"created_at"`
	CopiedAt   int64           `db:"copied_at"`
	Type       int16           `db:"type"`
	Content    json.RawMessage `db:"content"`
	Cover      string          `db:"cover"`
	Checked    bool            `db:"checked"`
	CategoryID *string         `db:"category_id"`
	TagText    *string         `db:"tag_text"`
	TagColor   *string         `db:"tag_color"`
}

func (r *NoteRepository) GetNotes(ctx context.Context, userID int64, offset int, filters []helpers.Filter) ([]entities.Note, error) {
	notes := []entities.Note{}

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	tags := []string{}
	category := ""
	for _, filter := range filters {
		if filter.Field == helpers.FilterKeyTag {
			tags = filter.Value.([]string)
		} else if filter.Field == helpers.FilterKeyCategory {
			category = filter.Value.(string)
		}
	}
	fmt.Println("Filters: ", filters)
	fmt.Println("Tags repo: ", tags)

	// Query for postgres
	filteredNotes := sq.Select("n.note_id", "n.creator_id", "n.title", "n.source", "n.original", "n.icon",
		"n.created_at", "n.copied_at", "n.type", "n.content", "n.cover", "n.checked", "n.category_id").
		From("notes n").
		Where(squirrel.Eq{"n.creator_id": userID})
	if category != "" {
		filteredNotes = filteredNotes.Where(squirrel.Eq{"n.category_id": category})
	}
	filteredNotes = filteredNotes.OrderBy("n.created_at DESC").
		Limit(10).
		Offset(uint64(offset))

	if len(tags) > 0 {
		filteredNotes = filteredNotes.
			Join("note_tag nt ON n.note_id = nt.note_id").
			Where(squirrel.Eq{"nt.tag_text": tags})
	}

	// Outer query to fetch all associated tags for the filtered notes
	query := sq.Select("fn.note_id", "fn.creator_id", "fn.title", "fn.source", "fn.original",
		"fn.icon", "fn.created_at", "fn.copied_at", "fn.type", "fn.content", "fn.cover",
		"fn.checked", "fn.category_id", "t.tag_text", "t.tag_color").
		FromSelect(filteredNotes, "fn").
		LeftJoin("note_tag nt ON fn.note_id = nt.note_id").
		LeftJoin("tags t ON nt.tag_text = t.tag_text AND t.owner_id = fn.creator_id")

	sql, args, err := query.ToSql()

	fmt.Println(sql, args)

	if err != nil {
		slog.Error("failed to build query: " + err.Error())
		return nil, err
	}

	notesInternal := []GetNotesInternal{}
	if err := r.DB.Select(&notesInternal, sql, args...); err != nil {
		return nil, err
	}

	if len(notesInternal) == 0 {
		return nil, nil
	}

	notes = append(notes, entities.Note{
		ID:        notesInternal[0].NoteID,
		CreatorID: notesInternal[0].CreatorID,
		Title:     notesInternal[0].Title,
		Source:    notesInternal[0].Source,
		Original:  notesInternal[0].Original,
		Icon:      notesInternal[0].Icon,
		CreatedAt: notesInternal[0].CreatedAt,
		CopiedAt:  notesInternal[0].CopiedAt,
		Type:      notesInternal[0].Type,
		Content:   notesInternal[0].Content,
		Cover:     notesInternal[0].Cover,
		Checked:   notesInternal[0].Checked,
	})
	if notesInternal[0].TagText != nil {
		notes[0].Tags = append(notes[0].Tags, entities.Tag{
			Text:  *notesInternal[0].TagText,
			Color: *notesInternal[0].TagColor,
		})
	}

	for _, noteInternal := range notesInternal[1:] {
		if noteInternal.TagText != nil && notes[len(notes)-1].ID == noteInternal.NoteID {
			notes[len(notes)-1].Tags = append(notes[len(notes)-1].Tags, entities.Tag{
				Text:  *noteInternal.TagText,
				Color: *noteInternal.TagColor,
			})
		} else {
			notes = append(notes, entities.Note{
				ID:         noteInternal.NoteID,
				CreatorID:  noteInternal.CreatorID,
				Title:      noteInternal.Title,
				Source:     noteInternal.Source,
				Original:   noteInternal.Original,
				Icon:       noteInternal.Icon,
				CreatedAt:  noteInternal.CreatedAt,
				CopiedAt:   noteInternal.CopiedAt,
				Type:       noteInternal.Type,
				Content:    noteInternal.Content,
				Cover:      noteInternal.Cover,
				Checked:    noteInternal.Checked,
				CategoryId: noteInternal.CategoryID,
			})
			if noteInternal.TagText != nil {
				notes[len(notes)-1].Tags = append(notes[len(notes)-1].Tags, entities.Tag{
					Text:  *noteInternal.TagText,
					Color: *noteInternal.TagColor,
				})
			}
		}
	}

	return notes, nil
}

func (r *NoteRepository) DeleteNotes(ctx context.Context, userID int64, noteIDs []string) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}
	for _, noteID := range noteIDs {
		if _, err := tx.Exec("DELETE FROM notes WHERE creator_id = $1 AND note_id = $2", userID, noteID); err != nil {
			return err
		}
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (r *NoteRepository) GetTags(ctx context.Context, userID int64) ([]entities.Tag, error) {
	tags := []entities.Tag{}
	if err := r.DB.Select(&tags, "SELECT * FROM tags WHERE owner_id = $1", userID); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *NoteRepository) CreateTag(ctx context.Context, tag *entities.Tag) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("INSERT INTO tags (owner_id, tag_text, tag_color) VALUES ($1, $2, $3)", tag.Owner, tag.Text, tag.Color); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (r *NoteRepository) DeleteTag(ctx context.Context, tag *entities.Tag) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("DELETE FROM tags WHERE creator_id = $1 AND tag_text = $2", tag.Owner, tag.Text); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (r *NoteRepository) RemoveTagFromNote(ctx context.Context, tag *entities.Tag, noteID string) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("DELETE FROM note_tag WHERE note_id = $1 AND tag_text = $2", noteID, tag.Text); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (r *NoteRepository) AddTagToNote(ctx context.Context, tag *entities.Tag, noteID string) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	fmt.Println("noteID: ", noteID)
	fmt.Printf("tag: %v\n", tag)

	if _, err := tx.Exec("INSERT INTO note_tag (note_id, tag_text, owner_id) VALUES ($1, $2, $3)", noteID, tag.Text, tag.Owner); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (r *NoteRepository) SetNoteCategory(ctx context.Context, userID int64, noteID string, categoryID string) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("UPDATE notes SET category_id = $1 WHERE note_id = $2 AND creator_id = $3", categoryID, noteID, userID); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
