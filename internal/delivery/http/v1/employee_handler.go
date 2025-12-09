package v1

import (
	"github.com/mohamedfawas/employee_management_system/internal/usecase"
)

type EmployeeHandler struct {
	employeeUsecase usecase.EmployeeUsecase
}

func NewEmployeeHandler(employeeUsecase usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{employeeUsecase: employeeUsecase}
}
