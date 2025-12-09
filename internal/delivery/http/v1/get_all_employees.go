package v1

import (
	"log"

	"github.com/labstack/echo/v4"
	apiresponse "github.com/mohamedfawas/employee_management_system/pkg/apiresponse"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

// GetAllEmployees retrieves all employees
// @Summary Get all employees
// @Description Retrieve a list of all employees in the system
// @Tags employees
// @Produce json
// @Success 200 {object} GetAllEmployeesResponseWrapper
// @Failure 500 {object} apiresponse.StandardResponse
// @Router /employees [get]
func (h *EmployeeHandler) GetAllEmployees(c echo.Context) error {
	employees, err := h.employeeUsecase.GetAllEmployees(c.Request().Context())
	if err != nil {
		if appError.ShouldLogError(err) {
			log.Printf("Error getting all employees: %v", err)
		}
		return apiresponse.Error(c, err, nil)
	}
	if len(employees) == 0 {
		return apiresponse.Success(c, "No employees found", nil)
	}
	employeesResponse := []GetAllEmployeesResponse{}
	for _, employee := range employees {
		employeesResponse = append(employeesResponse, GetAllEmployeesResponse{
			ID:        employee.ID,
			Name:      employee.Name,
			Position:  employee.Position,
			Salary:    employee.Salary,
			HiredDate: employee.HiredDate.Format("2006-01-02"),
			CreatedAt: employee.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return apiresponse.Success(c, "Employees retrieved successfully", employeesResponse)
}
