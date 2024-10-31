package domain

import (
	"errors"
	"time"
)

// Status defines the possible statuses for News
type NewsStatus string

const (
	Draft     NewsStatus = "draft"
	Deleted   NewsStatus = "deleted"
	Published NewsStatus = "published"
)

// Validate validates if the status is one of the predefined values
func (s NewsStatus) Validate() error {
	switch s {
	case Draft, Published, Deleted:
		return nil
	default:
		return errors.New("invalid status value")
	}
}

// News is representing the News data struct
type News struct {
	ID        int64       `json:"id"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Author    AuthorNews  `json:"author"` // just a little improvisation :)
	Status    NewsStatus  `json:"status"`
	UpdatedAt time.Time   `json:"updated_at"`
	CreatedAt time.Time   `json:"created_at"`
	Topics    []TopicNews `json:"topics"`
}

type NewsFilter struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	AuthorID  int64     `json:"author_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Limit     int64     `json:"limit"`
	Page      int64     `json:"page"`
	SortBy    string    `json:"sort_by"`    // e.g., "created_at"
	SortOrder string    `json:"sort_order"` // e.g., "asc" or "desc"
}
