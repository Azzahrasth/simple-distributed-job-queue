package entity

import "github.com/google/uuid"

// Job Status
const (
	PENDING   = "PENDING"
	RUNNING   = "RUNNING"
	COMPLETED = "COMPLETED"
	FAILED    = "FAILED"
	CANCELED  = "CANCELED"
)

type Job struct {
	ID       string `json:"id"`
	Task     string `json:"task"`
	Status   string `json:"status"`
	Attempts int32  `json:"attempts"`
}

type JobStatus struct {
	Pending   int32 `json:"pending"`
	Running   int32 `json:"running"`
	Failed    int32 `json:"failed"`
	Completed int32 `json:"completed"`
}

// NewJob creates a new Job instance
func NewJob(taskName string) *Job {
	return &Job{
		ID:       uuid.NewString(),
		Task:     taskName,
		Status:   PENDING,
		Attempts: 0,
	}
}
