package main

import (
	"time"
)

type ResponseNews struct {
	Data []NewsArticle  `json:"data"` // List of news articles
	Meta PaginationMeta `json:"meta"` // Metadata about the response
}

type ResponseTopic struct {
	Data []Topic        `json:"data"` // List of news articles
	Meta PaginationMeta `json:"meta"` // Metadata about the response
}

type NewsRequest struct {
	Title    string  `json:"title" validate:"required"`
	Content  string  `json:"content" validate:"required"`
	AuthorID int64   `json:"author_id" validate:"required"`
	Status   string  `json:"status" validate:"required"`
	TopicIDs []int64 `json:"topic_ids" validate:"required"`
}

type TopicRequest struct {
	Name string `json:"name" validate:"required"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type NewsArticle struct {
	ID        int       `json:"id"`         // ID of the news article
	Title     string    `json:"title"`      // Title of the news article
	Content   string    `json:"content"`    // Content of the news article (HTML)
	Author    Author    `json:"author"`     // Author of the news article
	Status    string    `json:"status"`     // Status of the news article
	UpdatedAt time.Time `json:"updated_at"` // Updated timestamp of the news article
	CreatedAt time.Time `json:"created_at"` // Created timestamp of the news article
	Topics    []Topic   `json:"topics"`     // List of topics associated with the news article
}

type Author struct {
	ID   int    `json:"id"`   // ID of the author
	Name string `json:"name"` // Name of the author
}

type Topic struct {
	ID   int    `json:"id"`   // ID of the topic
	Name string `json:"name"` // Name of the topic
}

type PaginationMeta struct {
	CurrentPage int64 `json:"current_page"` // Current page number
	TotalPages  int64 `json:"total_pages"`  // Total number of pages
	TotalData   int64 `json:"total_data"`   // Total number of data entries
}

// FetchNews handles GET requests to fetch news with optional filters
// @Summary Fetch news articles
// @Description Fetch news articles with optional filters like ID, title, status, author_id, start_date, end_date, sort_by, and sort_order.
// @Tags News
// @Accept json
// @Produce json
// @Param limit query int false "Limit the number of results"
// @Param page query int false "Page number for pagination"
// @Param id query int false "Filter by news ID"
// @Param title query string false "Filter by news title"
// @Param status query string false "Filter by news status"
// @Param author_id query int false "Filter by author ID"
// @Param start_date query string false "Filter news starting from this date (RFC3339 format)"
// @Param end_date query string false "Filter news until this date (RFC3339 format)"
// @Param sort_by query string false "Field to sort by"
// @Param sort_order query string false "Order of sorting (asc or desc)"
// @Success 200 {object} ResponseNews "Successful response with the list of news"
// @Router /news [get]
func FetchNews() {
}

// GetNewsByID handles GET requests to fetch a news article by its ID
// @Summary Get news article by ID
// @Description Fetch a news article by its unique ID
// @Tags News
// @Accept json
// @Produce json
// @Param id path int64 true "ID of the news article"
// @Success 200 {object} NewsArticle "Successful response with the news article"
// @Router /news/{id} [get]
func GetNewsByID() {}

// StoreNews handles POST requests to create a new news article
// @Summary Create a new news article
// @Description Create a news article with the provided details
// @Tags News
// @Accept json
// @Produce json
// @Param data body NewsRequest true "Create news request"
// @Success 201 {object} ResponseMessage "Success message"
// @Router /news [post]
func StoreNews() {}

// UpdateNews handles PUT requests to update an existing news article
// @Summary Update an existing news article
// @Description Update the news article with the provided ID using the details provided in the request body
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "News ID"
// @Param data body NewsRequest true "Update news request"
// @Success 200 {object} ResponseMessage "Success message"
// @Router /news/{id} [put]
func UpdateNews() {}

// DeleteNews handles DELETE requests to remove a news article by ID
// @Summary Delete a news article
// @Description Remove the news article with the specified ID
// @Tags News
// @Param id path int true "News ID"
// @Success 204 "No Content"
// @Router /news/{id} [delete]
func DeleteNews() {}

// FetchTopic handles GET requests to fetch topics with optional filters
// @Summary Fetch topics
// @Description Retrieve a list of topics with optional pagination and filtering
// @Tags Topic
// @Param limit query int false "Limit the number of topics returned" default(10)
// @Param page query int false "Page number for pagination" default(1)
// @Param id query int false "Filter by topic ID"
// @Param name query string false "Filter by topic name"
// @Param sort_by query string false "Field to sort by"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Success 200 {object} ResponseTopic "Successful response"
// @Router /topic [get]
func FetchTopic() {}

// GetTopicByID handles GET requests to fetch a topic by its ID
// @Summary Get topic by ID
// @Description Retrieve a topic by its ID
// @Tags Topic
// @Param id path int true "Topic ID"
// @Success 200 {object} Topic "Successful response with the topic"
// @Router /topic/{id} [get]
func GetTopicByID() {}

// StoreTopic handles POST requests to create a new topic
// @Summary Create a new topic
// @Description Add a new topic to the system
// @Tags Topic
// @Accept json
// @Produce json
// @Param topic body TopicRequest true "Create Topic Request"
// @Success 201 {object} ResponseMessage "Successful response message"
// @Router /topic [post]
func StoreTopic() {}

// UpdateTopic handles PUT requests to update an existing topic
// @Summary Update an existing topic
// @Description Update the details of a topic by its ID
// @Tags Topic
// @Accept json
// @Produce json
// @Param id path int true "Topic ID"
// @Param topic body TopicRequest true "Update Topic Request"
// @Success 200 {object} ResponseMessage "Successful response message"
// @Router /topic/{id} [put]
func UpdateTopic() {}

// DeleteTopic handles DELETE requests to remove a topic
// @Summary Delete a topic by ID
// @Description Delete an existing topic by its ID
// @Tags Topic
// @Param id path int true "Topic ID"
// @Success 204 "No Content"
// @Router /topic/{id} [delete]
func DeleteTopic() {}
