package documents

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted Status = "completed"
	StatusFailed Status = "failed"
)

type Document struct {
	ID 			string 		`json:"id"`
	Filename 	string 		`json:"filename"`
	ContentType string 		`json:"content_type"`
	SizeBytes 	int64 		`json:"size_bytes"`
	Status 		Status 		`json:"status"`
	LastError 	*string		`json:"last_error,omitempty"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInput struct {
	Filename    string `json:"filename"`
    ContentType string `json:"content_type"`
    SizeBytes   int64  `json:"size_bytes"`
}