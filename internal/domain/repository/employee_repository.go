package repository

import (
	"context"

	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employee *entity.Employee) error
	GetEmployeeById(ctx context.Context, id int) (*entity.Employee, error)
	GetAllEmployees(ctx context.Context) ([]*entity.Employee, error)
	UpdateEmployee(ctx context.Context, employee *entity.Employee) error
	DeleteEmployee(ctx context.Context, id int) error
}
