package usecase

import (
	"context"

	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

func (u *employeeUsecaseImpl) UpdateEmployee(ctx context.Context, employee *entity.Employee) (*entity.Employee, error) {
	if employee.ID <= 0 {
		return nil, appError.ErrInvalidEmployeeId
	}

	if employee.Name == "" || len(employee.Name) < 3 {
		return nil, appError.ErrInvalidName
	}
	if employee.Position == "" || len(employee.Position) < 3 {
		return nil, appError.ErrInvalidPosition
	}
	if employee.HiredDate.IsZero() {
		return nil, appError.ErrInvalidHiredDate
	}
	if employee.Salary <= 0 {
		return nil, appError.ErrInvalidSalary
	}

	updatedEmployee, err := u.employeeRepository.UpdateEmployee(ctx, employee)
	if updatedEmployee == nil {
		return nil, appError.ErrEmployeeNotFound
	}
	if err != nil {
		return nil, err
	}
	return updatedEmployee, nil
}
