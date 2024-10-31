package domain

// NewsTopic representing the NewsTopic relation data struct
type NewsTopic struct {
	NewsID  int64 `json:"news_id"`
	TopicID int64 `json:"topic_id"`
}
