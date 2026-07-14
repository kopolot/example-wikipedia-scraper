package api

import (
	types "example-wikipedia-scraper/internal/types/api"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateAndUnmarshalDTO[T any](body string) (*types.ApiResponse, *T) {
	var dto T
	if err := json.Unmarshal([]byte(body), &dto); err != nil {
		return BadRequestResponse(), nil
	}
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		var errs []types.ErrorDetail
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errs = ValidationErrorsToDetails(verrs)
		} else {
			errs = []types.ErrorDetail{{Message: err.Error()}}
		}
		return &types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Body: &types.ApiResponseBody{
				Success: false,
				Errors:  errs,
			},
		}, nil
	}
	return nil, &dto
}

func InternalErrorResponse() *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusInternalServerError,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: "Internal Server Error"}},
		},
	}
}

func NotFoundResponse(message string) *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusNotFound,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: message}},
		},
	}
}

func BadRequestResponseWithMsg(msg string) *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusBadRequest,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: msg}},
		},
	}
}

func BadRequestResponse() *types.ApiResponse {
	return BadRequestResponseWithMsg("Invalid request data")
}

func UnauthorizedResponse(msg string) *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusUnauthorized,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: msg}},
		},
	}
}

func ForbiddenResponseWithMsg(msg string) *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusForbidden,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: msg}},
		},
	}
}

func ForbiddenResponse() *types.ApiResponse {
	return ForbiddenResponseWithMsg("Forbidden")
}

func ConflictResponseWithMsg(msg string) *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusConflict,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: msg}},
		},
	}
}

func OkResponse(data any) *types.ApiResponse {
	return SuccesfulResponse(http.StatusOK, data)
}

func CreatedResponse(data any) *types.ApiResponse {
	return SuccesfulResponse(http.StatusCreated, data)
}

func SuccesfulResponse(status int, data any) *types.ApiResponse {
	if data == nil {
		data = map[string]any{}
	}
	return &types.ApiResponse{
		StatusCode: status,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    data,
		},
	}
}

func ValidationErrorsToDetails(verrs validator.ValidationErrors) []types.ErrorDetail {
	out := make([]types.ErrorDetail, 0, len(verrs))
	for _, v := range verrs {
		field := strings.ToLower(v.Field())
		out = append(out, types.ErrorDetail{
			Field:   field,
			Tag:     v.ActualTag(),
			Message: field + ": " + v.ActualTag(),
			Param:   v.Param(),
		})
	}
	return out
}

func getHeadersFromGinContext(c *gin.Context) map[string]string {
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	return headers
}

func getQueryParamsFromGinContext(c *gin.Context) map[string][]string {
	queryParams := make(map[string][]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values
		}
	}
	return queryParams
}
