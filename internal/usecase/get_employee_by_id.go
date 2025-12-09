package usecase

import (
	"context"

	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

func (u *employeeUsecaseImpl) GetEmployeeById(ctx context.Context, id int) (*entity.Employee, error) {
	if id <= 0 {
		return nil, appError.ErrInvalidEmployeeId
	}
	employee, err := u.employeeRepository.GetEmployeeById(ctx, id)
	if employee == nil {
		return nil, appError.ErrEmployeeNotFound
	}
	if err != nil {
		return nil, err
	}
	return employee, nil
}
