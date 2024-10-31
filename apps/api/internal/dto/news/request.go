package news

import (
	"github.com/bxcodec/go-clean-arch/domain"
	"time"
)

type CreateNewsReq struct {
	ID        int64             `json:"id"`
	Title     string            `json:"title" validate:"required"`
	Content   string            `json:"content" validate:"required"`
	AuthorID  int64             `json:"author_id" validate:"required"`
	Status    domain.NewsStatus `json:"status" validate:"required"`
	TopicIDs  []int64           `json:"topic_ids" validate:"required"`
	UpdatedAt time.Time         `json:"updated_at"`
	CreatedAt time.Time         `json:"created_at"`
}

type UpdateNewsReq struct {
	ID        *int64             `json:"id"`         // Pointer to allow for optional ID
	Title     *string            `json:"title"`      // Pointer to allow for optional title
	Content   *string            `json:"content"`    // Pointer to allow for optional content
	AuthorID  *int64             `json:"author_id"`  // Pointer to allow for optional author ID
	Status    *domain.NewsStatus `json:"status"`     // Pointer to allow for optional status
	TopicIDs  *[]int64           `json:"topic_ids"`  // Pointer to allow for optional topic IDs
	UpdatedAt *time.Time         `json:"updated_at"` // Pointer to allow for optional update timestamp
}
