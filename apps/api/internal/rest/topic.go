package rest

import (
	"context"
	"github.com/bxcodec/go-clean-arch/internal/dto"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"

	"github.com/bxcodec/go-clean-arch/domain"
)

// TopicService represent the news's use cases
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

// TopicHandler  represent the http handler for news
type TopicHandler struct {
	Service TopicService
}

// NewTopicHandler will initialize the news/ resources endpoint
func NewTopicHandler(e *echo.Echo, svc TopicService) {
	handler := &TopicHandler{
		Service: svc,
	}
	e.GET("/topic", handler.Fetch)
	e.POST("/topic", handler.Store)
	e.GET("/topic/:id", handler.GetByID)
	e.PUT("/topic/:id", handler.Update)
	e.DELETE("/topic/:id", handler.Delete)
}

func (a *TopicHandler) Fetch(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = defaultPage
	}

	idStr := c.QueryParam("id")
	name := c.QueryParam("name")
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")

	// Build TopicFilter
	filter := domain.TopicFilter{
		Limit: int64(limit),
		Page:  int64(page),
	}

	// Set optional filters if available
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

	// Context creation and service call
	ctx := c.Request().Context()
	listAr, totalData, err := a.Service.Fetch(ctx, filter)
	if err != nil {
		return c.JSON(dto.GetStatusCode(err), dto.ResponseError{Message: err.Error()})
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
func (a *TopicHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	listAr, _, err := a.Service.Fetch(ctx, domain.TopicFilter{ID: id, Page: 1, Limit: 1})
	if err != nil {
		return c.JSON(dto.GetStatusCode(err), dto.ResponseError{Message: err.Error()})
	}

	if len(listAr) == 0 {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	return c.JSON(http.StatusOK, listAr[0])
}

func isRequestTopicValid(m *domain.Topic) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the news by given request body
func (a *TopicHandler) Store(c echo.Context) (err error) {
	var createTopicReq domain.Topic
	err = c.Bind(&createTopicReq)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if _, err = isRequestTopicValid(&createTopicReq); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.Service.Store(ctx, &createTopicReq)
	if err != nil {
		return c.JSON(dto.GetStatusCode(err), dto.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success create topic"})
}

// Update will update news by given param
func (a *TopicHandler) Update(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)

	updateTopicReq := domain.Topic{
		ID: id,
	}
	err = c.Bind(&updateTopicReq)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	err = a.Service.Update(ctx, &updateTopicReq)
	if err != nil {
		return c.JSON(dto.GetStatusCode(err), dto.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success update topic"})
}

// Delete will delete news by given param
func (a *TopicHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = a.Service.Delete(ctx, id)
	if err != nil {
		return c.JSON(dto.GetStatusCode(err), dto.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
