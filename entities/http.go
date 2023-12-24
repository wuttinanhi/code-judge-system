package entities

type HttpError struct {
	Message string `json:"message"`
}

type PaginationOptions struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginationResult[T any] struct {
	Total int `json:"total"`
	Items []T `json:"items"`
}
