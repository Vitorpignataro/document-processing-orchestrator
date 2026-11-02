package documents

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Vitorpignataro/document-processing-orchestrator/internal/jobs"
	"github.com/google/uuid"
)

var ErrInvalidInput =  errors.New("invalid document input")

type Service struct{
	db *sql.DB
	repo *Repository
}

func NewService(
	db *sql.DB,
	repo *Repository,
) *Service{
	return &Service{
		db: db,
		repo: repo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	input CreateInput,
) (*Document, error) {
	input.Filename = strings.TrimSpace(input.Filename)
	input.ContentType = strings.TrimSpace(input.ContentType)

	if input.Filename == "" ||
        input.ContentType == "" ||
        input.SizeBytes <= 0 ||
        len(input.Filename) > 255 ||
        len(input.ContentType) > 100 {
        return nil, ErrInvalidInput
    }

	documentID := uuid.NewString()
	jobID := uuid.NewString()

	now := time.Now().UTC()

	payload, err := json.Marshal(
		jobs.ValidateDocumentPayload{
			DocumentID: documentID,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERTO INTO documents (
			id,
			filename,
			content_type,
			size_bytes,
			status,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ? , 'pending', ?, ?)
	`,
		documentID,
		input.Filename,
		input.ContentType,
		input.SizeBytes,
		now,
		now,
	)

	if err != nil {
		return nil, fmt.Errorf("insert document: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO jobs (
			id,
			type,
			payload,
			status,
            attempts,
			max_attempts,
            available_at,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, 'pending', 0, 3, ?, ?, ?)
	`,
		jobID,
		jobs.TypeValidateDocument,
		payload,
		now,
		now,
		now,
	)

	if err != nil {
		return nil, fmt.Errorf("insert job: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return s.repo.GetByID(ctx, documentID)
}


