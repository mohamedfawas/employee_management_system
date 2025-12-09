package v1

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	apiresponse "github.com/mohamedfawas/employee_management_system/pkg/apiresponse"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

// DeleteEmployee deletes an employee by ID
// @Summary Delete an employee
// @Description Remove an employee record by its ID
// @Tags employees
// @Param id path int true "Employee ID"
// @Produce json
// @Success 204 "No Content"
// @Failure 400 {object} apiresponse.StandardResponse
// @Failure 404 {object} apiresponse.StandardResponse
// @Failure 500 {object} apiresponse.StandardResponse
// @Router /employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c echo.Context) error {
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

	err = h.employeeUsecase.DeleteEmployee(c.Request().Context(), id)
	if err != nil {
		if appError.ShouldLogError(err) {
			log.Printf("Error deleting employee: %v", err)
		}
		return apiresponse.Error(c, err, nil)
	}

	return apiresponse.DeletedResource(c, "Employee deleted successfully")
}
