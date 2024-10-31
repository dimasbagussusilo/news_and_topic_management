package dto

type PaginationMeta struct {
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
	TotalData   int64 `json:"total_data"`
}

type Response struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
