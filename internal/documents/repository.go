package documents

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("document not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*Document, error) {
	const query = `
	SELECT
		id,
		filename,
		content_type,
		size_bytes,
		status,
		last_error,
		created_at,
		updated_at,
	FROM documents
	WHERE id = ?
	`
	var document Document

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&document.ID,
		&document.Filename,
        &document.ContentType,
        &document.SizeBytes,
        &document.Status,
        &document.LastError,
        &document.CreatedAt,
        &document.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}

	return &document, nil
}