package api

type ApiResponse struct {
	StatusCode int
	Body       *ApiResponseBody
	Headers    map[string]string
}

type ApiResponseBody struct {
	Success bool          `json:"success"`
	Data    any           `json:"data,omitempty"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Param   string `json:"param,omitempty"`
}
