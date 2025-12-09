package usecase

import (
	"context"

	domaincache "github.com/mohamedfawas/employee_management_system/internal/domain/cache"
	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
	"github.com/mohamedfawas/employee_management_system/internal/domain/repository"
)

type EmployeeUsecase interface {
	CreateEmployee(ctx context.Context, employee *entity.Employee) (*entity.Employee, error)
	GetEmployeeById(ctx context.Context, id int) (*entity.Employee, error)
	GetAllEmployees(ctx context.Context) ([]*entity.Employee, error)
	UpdateEmployee(ctx context.Context, employee *entity.Employee) (*entity.Employee, error)
	DeleteEmployee(ctx context.Context, id int) error
}

type employeeUsecaseImpl struct {
	employeeRepository repository.EmployeeRepository
	cache              domaincache.Cache
}

func NewEmployeeUsecase(employeeRepository repository.EmployeeRepository, cache domaincache.Cache) EmployeeUsecase {
	return &employeeUsecaseImpl{
		employeeRepository: employeeRepository,
		cache:              cache,
	}
}
