package v1

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	apiresponse "github.com/mohamedfawas/employee_management_system/pkg/apiresponse"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

// GetEmployeeById retrieves an employee by ID
// @Summary Get employee by ID
// @Description Fetch a single employee using the ID provided in the URL path
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} GetEmployeeByIdResponseWrapper
// @Failure 400 {object} apiresponse.StandardResponse
// @Failure 404 {object} apiresponse.StandardResponse
// @Failure 500 {object} apiresponse.StandardResponse
// @Router /employees/{id} [get]
func (h *EmployeeHandler) GetEmployeeById(c echo.Context) error {
	idParam := c.Param("id")
	if idParam == "" {
		return apiresponse.Error(c,
			appError.ErrMissingRequiredFields,
			map[string]string{
				"id": "ID is required in the URL path",
			})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return apiresponse.Error(c,
			appError.ErrInvalidEmployeeId,
			map[string]string{
				"id": "ID must be a valid number",
			})
	}

	employee, err := h.employeeUsecase.GetEmployeeById(c.Request().Context(), id)
	if err != nil {
		if appError.ShouldLogError(err) {
			log.Printf("Error getting employee by id: %v", err)
		}
		return apiresponse.Error(c, err, nil)
	}

	employeeResponse := GetEmployeeByIdResponse{
		ID:        employee.ID,
		Name:      employee.Name,
		Position:  employee.Position,
		Salary:    employee.Salary,
		HiredDate: employee.HiredDate.Format("2006-01-02"),
		CreatedAt: employee.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return apiresponse.Success(c, "Employee retrieved successfully", employeeResponse)
}
