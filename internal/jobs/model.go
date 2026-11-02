package jobs

import (
	"encoding/json"
	"time"
)

type Status string

const (
    StatusPending    Status = "pending"
    StatusProcessing Status = "processing"
    StatusCompleted  Status = "completed"
    StatusFailed     Status = "failed"
)

const TypeValidateDocument = "VALIDATE_DOCUMENT"

type Job struct {
	ID string
	Type string
	Payload json.RawMessage
	Status	Status
	Attempts int
	LeaseOwner string

	CreatedAt time.Time
}

type ValidateDocumentPayload struct {
	DocumentID string `json:"document_id"`
}