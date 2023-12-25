package entities

type PaginationOptions struct {
	Page  int    `json:"page" validate:"min=1"`
	Limit int    `json:"limit" validate:"min=1,max=100"`
	Order string `json:"order_by" validate:"oneof=asc desc"`
	Sort  string `json:"sort_by" validate:"string"`
}

type PaginationResult[T any] struct {
	Total int `json:"total"`
	Items []T `json:"items"`
}
