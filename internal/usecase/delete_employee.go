package usecase

import (
	"context"
	"errors"

	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

func (u *employeeUsecaseImpl) DeleteEmployee(ctx context.Context, id int) error {
	if id <= 0 {
		return appError.ErrInvalidEmployeeId
	}
	err := u.employeeRepository.DeleteEmployee(ctx, id)
	if err != nil {
		if errors.Is(err, appError.ErrEmployeeNotFound) {
			return appError.ErrEmployeeNotFound
		}
		return err
	}
	return nil
}
