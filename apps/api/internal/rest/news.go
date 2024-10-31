package rest

import (
	"context"
	"errors"
	"net/http"
	"strconv"

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
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.News, string, error)
	GetByID(ctx context.Context, id int64) (domain.News, error)
	Update(ctx context.Context, ar *domain.News) error
	GetByTitle(ctx context.Context, title string) (domain.News, error)
	Store(context.Context, *domain.News) error
	Delete(ctx context.Context, id int64) error
}

// NewsHandler  represent the http handler for news
type NewsHandler struct {
	Service NewsService
}

const defaultNum = 10

// NewNewsHandler will initialize the news/ resources endpoint
func NewNewsHandler(e *echo.Echo, svc NewsService) {
	handler := &NewsHandler{
		Service: svc,
	}
	e.GET("/news", handler.FetchNews)
	e.POST("/news", handler.Store)
	e.GET("/news/:id", handler.GetByID)
	e.DELETE("/news/:id", handler.Delete)
}

// FetchNews will fetch the news based on given params
func (a *NewsHandler) FetchNews(c echo.Context) error {

	numS := c.QueryParam("num")
	num, err := strconv.Atoi(numS)
	if err != nil || num == 0 {
		num = defaultNum
	}

	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := a.Service.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

// GetByID will get news by given id
func (a *NewsHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	art, err := a.Service.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *domain.News) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the news by given request body
func (a *NewsHandler) Store(c echo.Context) (err error) {
	var news domain.News
	err = c.Bind(&news)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if _, err = isRequestValid(&news); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.Service.Store(ctx, &news)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, news)
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
