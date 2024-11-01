package domain

import "time"

// Topic representing the Topic data struct
type Topic struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TopicNews representing the TopicNews data struct
type TopicNews struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TopicFilter struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Limit     int64  `json:"limit"`
	Page      int64  `json:"page"`
	SortBy    string `json:"sort_by"`    // e.g., "created_at"
	SortOrder string `json:"sort_order"` // e.g., "asc" or "desc"
}
