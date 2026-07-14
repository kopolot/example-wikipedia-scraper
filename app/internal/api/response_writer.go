package api

import (
	types "example-wikipedia-scraper/internal/types/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeApiResponse(c *gin.Context, response *types.ApiResponse) {
	if response == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for key, value := range response.Headers {
		c.Header(key, value)
	}

	if response.StatusCode >= http.StatusMultipleChoices && response.StatusCode < http.StatusBadRequest {
		if location := response.Headers["Location"]; location != "" {
			c.Redirect(response.StatusCode, location)
			return
		}
	}

	if response.Body == nil {
		c.Status(response.StatusCode)
		return
	}

	c.JSON(response.StatusCode, response.Body)
}
