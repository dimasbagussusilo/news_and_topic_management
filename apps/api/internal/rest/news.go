package rest

import (
	"context"
	"errors"
	"github.com/bxcodec/go-clean-arch/internal/dto"
	"github.com/bxcodec/go-clean-arch/internal/dto/news"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"

	"github.com/bxcodec/go-clean-arch/domain"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// NewsService represent the news's use cases
//
//go:generate mockery --name NewsService
type NewsService interface {
	Fetch(ctx context.Context, filter domain.NewsFilter) ([]domain.News, int64, error)
	GetByID(ctx context.Context, id int64) (domain.News, error)
	Update(ctx context.Context, ar *news.UpdateNewsReq) error
	GetByTitle(ctx context.Context, title string) (domain.News, error)
	Store(context.Context, *news.CreateNewsReq) error
	Delete(ctx context.Context, id int64) error
}

// NewsHandler  represent the http handler for news
type NewsHandler struct {
	Service NewsService
}

const defaultLimit = 10
const defaultPage = 1

// NewNewsHandler will initialize the news/ resources endpoint
func NewNewsHandler(e *echo.Echo, svc NewsService) {
	handler := &NewsHandler{
		Service: svc,
	}
	e.GET("/news", handler.Fetch)
	e.POST("/news", handler.Store)
	e.GET("/news/:id", handler.GetByID)
	e.PUT("/news/:id", handler.Update)
	e.DELETE("/news/:id", handler.Delete)
}

func (a *NewsHandler) Fetch(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = defaultPage
	}

	idStr := c.QueryParam("id")
	title := c.QueryParam("title")
	status := c.QueryParam("status")
	authorIDStr := c.QueryParam("author_id")
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")

	// Build NewsFilter
	filter := domain.NewsFilter{
		Limit: int64(limit),
		Page:  int64(page),
	}

	// Set optional filters if available
	if idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			filter.ID = id
		}
	}
	if title != "" {
		filter.Title = title
	}
	if status != "" {
		filter.Status = status
	}
	if authorIDStr != "" {
		if authorID, err := strconv.ParseInt(authorIDStr, 10, 64); err == nil {
			filter.AuthorID = authorID
		}
	}
	if startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filter.StartDate = startDate
		}
	}
	if endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filter.EndDate = endDate
		}
	}
	if sortBy != "" {
		filter.SortBy = sortBy
	}
	if sortOrder != "" {
		filter.SortOrder = sortOrder
	}

	// Context creation and service call
	ctx := c.Request().Context()
	listAr, totalData, err := a.Service.Fetch(ctx, filter)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	// Calculate the total pages
	totalPages := (int64(totalData) + int64(filter.Limit) - 1) / int64(filter.Limit)

	// Construct the response
	response := dto.Response{
		Data: listAr,
		Meta: dto.PaginationMeta{
			CurrentPage: filter.Page,
			TotalPages:  totalPages,
			TotalData:   totalData,
		},
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID will get news by given id
func (a *NewsHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	listAr, _, err := a.Service.Fetch(ctx, domain.NewsFilter{ID: id, Page: 1, Limit: 1})
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	if len(listAr) == 0 {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	return c.JSON(http.StatusOK, listAr[0])
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
func (a *NewsHandler) Store(c echo.Context) (err error) {
	var createNewsReq news.CreateNewsReq
	err = c.Bind(&createNewsReq)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if _, err = isRequestValid(&createNewsReq); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.Service.Store(ctx, &createNewsReq)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success create news"})
}

// Update will update news by given param
func (a *NewsHandler) Update(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)

	updateNewsReq := news.UpdateNewsReq{
		ID: &id,
	}
	err = c.Bind(&updateNewsReq)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	err = a.Service.Update(ctx, &updateNewsReq)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success update news"})
}

// Delete will delete news by given param
func (a *NewsHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = a.Service.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch {
	case errors.Is(err, domain.ErrInternalServerError):
		return http.StatusInternalServerError
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
