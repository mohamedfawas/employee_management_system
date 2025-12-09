package apiresponse

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

// StandardResponse represents a common API response
// swagger:model StandardResponse
type StandardResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

type ErrorInfo struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func Success(c echo.Context, msg string, data interface{}) error {
	return c.JSON(http.StatusOK, StandardResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: getRequestID(c),
	})
}

func DeletedResource(c echo.Context, msg string) error {
	return c.JSON(http.StatusNoContent, StandardResponse{
		Success:   true,
		Message:   msg,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: getRequestID(c),
	})
}

func Error(c echo.Context, err error, details map[string]string) error {
	var appErr *appError.AppError
	if errors.As(err, &appErr) {
		return c.JSON(appErr.HTTPStatusCode, StandardResponse{
			Success: false,
			Message: "Request failed",
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.PublicMsg,
				Details: details,
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: getRequestID(c),
		})
	}

	return c.JSON(http.StatusInternalServerError, StandardResponse{
		Success: false,
		Message: "Request failed",
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Something went wrong. Please try again later.",
			Details: nil, // never leak internal error info
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: getRequestID(c),
	})
}

func getRequestID(c echo.Context) string {
	rid := c.Request().Header.Get("X-Request-ID")
	if rid == "" {
		return ""
	}
	return rid
}
