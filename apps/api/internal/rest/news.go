package rest

import (
	"context"
	"encoding/json"
	"github.com/bxcodec/go-clean-arch/internal/dto"
	"github.com/bxcodec/go-clean-arch/internal/dto/news"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
)

// NewsService represents the news's use cases
type NewsService interface {
	Fetch(ctx context.Context, filter domain.NewsFilter) ([]domain.News, int64, error)
	GetByID(ctx context.Context, id int64) (domain.News, error)
	Update(ctx context.Context, ar *news.UpdateNewsReq) error
	GetByTitle(ctx context.Context, title string) (domain.News, error)
	Store(context.Context, *news.CreateNewsReq) error
	Delete(ctx context.Context, id int64) error
}

// NewsHandler represents the HTTP handler for news
type NewsHandler struct {
	Service NewsService
}

const (
	defaultLimit = 10
	defaultPage  = 1
)

// NewNewsHandler initializes the news resources endpoints
func NewNewsHandler(mux *http.ServeMux, svc NewsService) {
	handler := &NewsHandler{
		Service: svc,
	}
	mux.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.Fetch(w, r)
		case http.MethodPost:
			handler.Store(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/news/", handler.NewsHandler) // Combines GetByID, Update, and Delete based on HTTP method
}

// Fetch handles GET requests to fetch news with optional filters
func (a *NewsHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = defaultPage
	}

	// Parse optional filters from query parameters
	filter := domain.NewsFilter{
		Limit: int64(limit),
		Page:  int64(page),
	}

	// Set optional filters
	if idStr := r.URL.Query().Get("id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			filter.ID = id
		}
	}
	if title := r.URL.Query().Get("title"); title != "" {
		filter.Title = title
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = status
	}
	if authorIDStr := r.URL.Query().Get("author_id"); authorIDStr != "" {
		if authorID, err := strconv.ParseInt(authorIDStr, 10, 64); err == nil {
			filter.AuthorID = authorID
		}
	}
	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filter.StartDate = startDate
		}
	}
	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filter.EndDate = endDate
		}
	}
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}
	if sortOrder := r.URL.Query().Get("sort_order"); sortOrder != "" {
		filter.SortOrder = sortOrder
	}

	// Fetch news using the service
	ctx := r.Context()
	listAr, totalData, err := a.Service.Fetch(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), dto.GetStatusCode(err))
		return
	}

	// Calculate total pages
	totalPages := ((totalData) + (filter.Limit) - 1) / (filter.Limit)

	// Construct response
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

// NewsHandler routes based on HTTP method for ID-based operations
func (a *NewsHandler) NewsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/news/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.GetByID(w, r, id)
	case http.MethodPut:
		a.Update(w, r, id)
	case http.MethodDelete:
		a.Delete(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetByID retrieves news by the given ID
func (a *NewsHandler) GetByID(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	newsItem, _, err := a.Service.Fetch(ctx, domain.NewsFilter{ID: id, Page: 1, Limit: 1})
	if err != nil || len(newsItem) == 0 {
		http.Error(w, "News not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newsItem[0])
	if err != nil {
		return
	}
}

func isRequestValid(m *news.CreateNewsReq) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the news by given request body
func (a *NewsHandler) Store(w http.ResponseWriter, r *http.Request) {
	var createNewsReq news.CreateNewsReq
	err := json.NewDecoder(r.Body).Decode(&createNewsReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if _, err = isRequestValid(&createNewsReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = a.Service.Store(ctx, &createNewsReq)
	if err != nil {
		http.Error(w, err.Error(), dto.GetStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "success create news"})
	if err != nil {
		return
	}
}

// Update updates the news based on the given request
func (a *NewsHandler) Update(w http.ResponseWriter, r *http.Request, id int64) {
	var updateNewsReq news.UpdateNewsReq
	updateNewsReq.ID = &id

	if err := json.NewDecoder(r.Body).Decode(&updateNewsReq); err != nil {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	ctx := r.Context()
	if err := a.Service.Update(ctx, &updateNewsReq); err != nil {
		http.Error(w, err.Error(), dto.GetStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"message": "success update news"})
	if err != nil {
		return
	}
}

// Delete removes the news by the given ID
func (a *NewsHandler) Delete(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	if err := a.Service.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), dto.GetStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
