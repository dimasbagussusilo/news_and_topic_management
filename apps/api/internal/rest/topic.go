package rest

import (
	"context"
	"encoding/json"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/dto"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

// TopicService represents the topic's use cases
//
//go:generate mockery --name TopicService
type TopicService interface {
	Fetch(ctx context.Context, filter domain.TopicFilter) ([]domain.Topic, int64, error)
	GetByID(ctx context.Context, id int64) (domain.Topic, error)
	Update(ctx context.Context, ar *domain.Topic) error
	GetByTitle(ctx context.Context, title string) (domain.Topic, error)
	Store(context.Context, *domain.Topic) error
	Delete(ctx context.Context, id int64) error
}

// TopicHandler represents the HTTP handler for topics
type TopicHandler struct {
	Service TopicService
}

// NewTopicHandler initializes the topic resources endpoints
func NewTopicHandler(mux *http.ServeMux, svc TopicService) {
	handler := &TopicHandler{
		Service: svc,
	}
	mux.HandleFunc("/topic", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.Fetch(w, r)
		case http.MethodPost:
			handler.Store(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/topic/", handler.HandleTopicByID)
}

// HandleTopicByID routes requests to the appropriate handler based on HTTP method
func (a *TopicHandler) HandleTopicByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.GetByID(w, r)
	case http.MethodPut:
		a.Update(w, r)
	case http.MethodDelete:
		a.Delete(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (a *TopicHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = defaultPage
	}

	idStr := query.Get("id")
	name := query.Get("name")
	sortBy := query.Get("sort_by")
	sortOrder := query.Get("sort_order")

	// Build TopicFilter
	filter := domain.TopicFilter{
		Limit: int64(limit),
		Page:  int64(page),
	}

	if idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			filter.ID = id
		}
	}
	if name != "" {
		filter.Name = name
	}
	if sortBy != "" {
		filter.SortBy = sortBy
	}
	if sortOrder != "" {
		filter.SortOrder = sortOrder
	}

	ctx := r.Context()
	listAr, totalData, err := a.Service.Fetch(ctx, filter)
	if err != nil {
		http.Error(w, dto.ResponseError{Message: err.Error()}.Message, dto.GetStatusCode(err))
		return
	}

	totalPages := (int64(totalData) + int64(filter.Limit) - 1) / int64(filter.Limit)

	response := dto.Response{
		Data: listAr,
		Meta: dto.PaginationMeta{
			CurrentPage: filter.Page,
			TotalPages:  totalPages,
			TotalData:   totalData,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (a *TopicHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/topic/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	ctx := r.Context()
	listAr, _, err := a.Service.Fetch(ctx, domain.TopicFilter{ID: id, Page: 1, Limit: 1})
	if err != nil || len(listAr) == 0 {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(listAr[0])
	if err != nil {
		return
	}
}

func isRequestTopicValid(m *domain.Topic) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	return err == nil, err
}

func (a *TopicHandler) Store(w http.ResponseWriter, r *http.Request) {
	var createTopicReq domain.Topic
	err := json.NewDecoder(r.Body).Decode(&createTopicReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if _, err = isRequestTopicValid(&createTopicReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = a.Service.Store(ctx, &createTopicReq)
	if err != nil {
		http.Error(w, dto.ResponseError{Message: err.Error()}.Message, dto.GetStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "success create topic"})
	if err != nil {
		return
	}
}

func (a *TopicHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/topic/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	updateTopicReq := domain.Topic{
		ID: id,
	}
	err = json.NewDecoder(r.Body).Decode(&updateTopicReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	ctx := r.Context()
	err = a.Service.Update(ctx, &updateTopicReq)
	if err != nil {
		http.Error(w, dto.ResponseError{Message: err.Error()}.Message, dto.GetStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "success update topic"})
	if err != nil {
		return
	}
}

func (a *TopicHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/topic/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	ctx := r.Context()
	err = a.Service.Delete(ctx, id)
	if err != nil {
		http.Error(w, dto.ResponseError{Message: err.Error()}.Message, dto.GetStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
