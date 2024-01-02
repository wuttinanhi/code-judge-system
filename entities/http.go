package entities

type HttpError struct {
	Message string `json:"message"`
}

type HttpBadRequest struct {
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}
