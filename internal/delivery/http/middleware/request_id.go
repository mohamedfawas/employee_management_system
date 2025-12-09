package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mohamedfawas/employee_management_system/pkg/constants"
)

func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			requestID := c.Request().Header.Get(constants.HeaderRequestID)

			if requestID == "" {
				requestID = uuid.New().String()
			}

			c.Set(constants.ContextKeyRequestID, requestID)

			c.Response().Header().Set(constants.HeaderRequestID, requestID)

			return next(c)
		}
	}
}
