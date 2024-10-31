package domain

import (
	"errors"
	"time"
)

// Status defines the possible statuses for News
type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived" // it's deleted in requirement
)

// Validate validates if the status is one of the predefined values
func (s Status) Validate() error {
	switch s {
	case StatusDraft, StatusPublished, StatusArchived:
		return nil
	default:
		return errors.New("invalid status value")
	}
}

// News is representing the News data struct
type News struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title" validate:"required"`
	Content   string     `json:"content" validate:"required"`
	Author    AuthorNews `json:"author"` // just a little improvisation :)
	Status    Status     `json:"status"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	Topics    []TopicNews
}
