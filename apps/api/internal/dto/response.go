package dto

import (
	"errors"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PaginationMeta struct {
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
	TotalData   int64 `json:"total_data"`
}

type Response struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type ResponseError struct {
	Message string `json:"message"`
}

func GetStatusCode(err error) int {
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
