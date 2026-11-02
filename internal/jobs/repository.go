package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoJobs = errors.New("no jobns available")
var ErrOwnershipLost = errors.New("job ownership lost")

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Acquire (
	ctx context.Context,
	workerID string,
) (*Job, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		 return nil, fmt.Errorf("begin acquire transaction: %w", err)
	}
	defer tx.Rollback()

	var job Job

	err = tx.QueryRowContext(ctx, `
        SELECT
            id, type, payload, status,
            attempts, created_at
        FROM jobs
        WHERE status = 'pending'
          AND available_at <= UTC_TIMESTAMP(6)
        ORDER BY available_at, created_at, id
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    `).Scan(
        &job.ID,
        &job.Type,
        &job.Payload,
        &job.Status,
        &job.Attempts,
        &job.CreatedAt,
    )

	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrNoJobs
	}

	if err != nil {
		return nil, fmt.Errorf("select job: %w", err)
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE jobs
		SET
			 status = 'processing',
            lease_owner = ?,
            lease_expires_at = DATE_ADD(
                UTC_TIMESTAMP(6),
                INTERVAL 5 MINUTE
            ),
            last_heartbeat_at = UTC_TIMESTAMP(6),
            attempts = attempts + 1,
            updated_at = UTC_TIMESTAMP(6)
        WHERE id = ?
          AND status = 'pending'
    `, workerID, job.ID)

	if err != nil {
        return nil, fmt.Errorf("acquire job: %w", err)
    }

	affected, err := result.RowsAffected()
    if err != nil {
        return nil, err
    }

	if affected != 1 {
        return nil, ErrOwnershipLost
    }

	 if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("commit acquisition: %w", err)
    }

    job.Status = StatusProcessing
    job.LeaseOwner = workerID
    job.Attempts++

    return &job, nil
}
