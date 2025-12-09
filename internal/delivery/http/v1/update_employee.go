package v1

import (
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
	apiresponse "github.com/mohamedfawas/employee_management_system/pkg/apiresponse"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

// UpdateEmployee updates an existing employee.
//
// @Summary Update an employee
// @Description Update employee details by ID
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param payload body UpdateEmployeeRequest true "Update employee payload"
// @Success 200 {object} UpdateEmployeeResponseWrapper
// @Failure 400 {object} apiresponse.StandardResponse
// @Failure 404 {object} apiresponse.StandardResponse
// @Failure 500 {object} apiresponse.StandardResponse
// @Router /employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c echo.Context) error {
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

	var req UpdateEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return apiresponse.Error(c,
			appError.ErrMissingRequiredFields,
			map[string]string{
				"name":       "Name is required and must be at least 3 characters long",
				"position":   "Position is required and must be at least 3 characters long",
				"salary":     "Salary is required and must be a valid number",
				"hired_date": "Hired date is required and must be a valid date",
			})
	}

	hiredDate, err := time.Parse("2006-01-02", req.HiredDate)
	if err != nil {
		return apiresponse.Error(c,
			appError.ErrInvalidHiredDate,
			map[string]string{
				"hired_date": "Date format is invalid , expected format: YYYY-MM-DD",
			},
		)
	}

	employee := &entity.Employee{
		ID:        id,
		Name:      req.Name,
		Position:  req.Position,
		Salary:    req.Salary,
		HiredDate: hiredDate,
	}

	updatedEmployee, err := h.employeeUsecase.UpdateEmployee(c.Request().Context(), employee)
	if err != nil {
		if appError.ShouldLogError(err) {
			log.Printf("Error updating employee: %v", err)
		}
		return apiresponse.Error(c, err, nil)
	}

	updatedEmployeeResponse := UpdateEmployeeResponse{
		ID:        updatedEmployee.ID,
		Name:      updatedEmployee.Name,
		Position:  updatedEmployee.Position,
		Salary:    updatedEmployee.Salary,
		HiredDate: updatedEmployee.HiredDate.Format("2006-01-02"),
		UpdatedAt: updatedEmployee.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return apiresponse.Success(c, "Employee updated successfully", updatedEmployeeResponse)
}
