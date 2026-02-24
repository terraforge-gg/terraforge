package dto

type ProblemDetails struct {
	Title  string            `json:"title"`
	Status int               `json:"status"`
	Detail string            `json:"detail"`
	Errors map[string]string `json:"errors,omitempty"`
}
