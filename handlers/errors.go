package handlers

import (
	"net/http"

	"agentcore/models"
	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, status int, code, message string, details interface{}) {
	envelope := models.ErrorEnvelope{
		Error: models.ErrorPayload{
			Code:    code,
			Message: message,
			Details: details,
		},
		RequestID: c.GetHeader("X-Request-ID"),
	}
	c.JSON(status, envelope)
}

func WriteValidationError(c *gin.Context, message string, details interface{}) {
	writeError(c, http.StatusBadRequest, "validation_error", message, details)
}

func WriteInternalError(c *gin.Context, message string) {
	writeError(c, http.StatusInternalServerError, "internal_error", message, nil)
}

func WriteNotFoundError(c *gin.Context, message string) {
	writeError(c, http.StatusNotFound, "not_found", message, nil)
}
