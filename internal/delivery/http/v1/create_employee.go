package v1

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
	apiresponse "github.com/mohamedfawas/employee_management_system/pkg/apiresponse"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Creates a new employee and stores it in the database.
// @Tags Employees
// @Accept json
// @Produce json
// @Param payload body CreateEmployeeRequest true "Employee create payload"
// @Success 200 {object} CreateEmployeeResponseWrapper
// @Failure 400 {object} apiresponse.StandardResponse
// @Failure 500 {object} apiresponse.StandardResponse
// @Router /employees [post]
func (h *EmployeeHandler) CreateEmployee(c echo.Context) error {
	var req CreateEmployeeRequest
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
		Name:      req.Name,
		Position:  req.Position,
		Salary:    req.Salary,
		HiredDate: hiredDate,
	}

	createdEmployee, err := h.employeeUsecase.CreateEmployee(c.Request().Context(), employee)
	if err != nil {
		if appError.ShouldLogError(err) {
			log.Printf("Error creating employee: %v", err)
		}
		return apiresponse.Error(c, err, nil)
	}

	createdEmployeeResponse := CreateEmployeeResponse{
		ID:        createdEmployee.ID,
		Name:      createdEmployee.Name,
		Position:  createdEmployee.Position,
		Salary:    createdEmployee.Salary,
		HiredDate: createdEmployee.HiredDate.Format("2006-01-02"),
		CreatedAt: createdEmployee.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return apiresponse.Success(c, "Employee created successfully", createdEmployeeResponse)
}
