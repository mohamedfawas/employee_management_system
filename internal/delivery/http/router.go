package http

import (
	"github.com/labstack/echo/v4"
	v1 "github.com/mohamedfawas/employee_management_system/internal/delivery/http/v1"
)

func RegisterRoutes(e *echo.Echo, h *v1.EmployeeHandler) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/employees", h.CreateEmployee)
		v1.GET("/employees/:id", h.GetEmployeeById)
		v1.GET("/employees", h.GetAllEmployees)
		v1.PUT("/employees/:id", h.UpdateEmployee)
		v1.DELETE("/employees/:id", h.DeleteEmployee)
	}
}
