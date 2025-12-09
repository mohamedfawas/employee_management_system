package apperror

import (
	"errors"
	"net/http"

	"github.com/mohamedfawas/employee_management_system/pkg/constants"
)

var (
	ErrMissingRequiredFields = &AppError{
		Err:            errors.New("missing required fields"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Missing required fields",
	}
	ErrInvalidName = &AppError{
		Err:            errors.New("invalid name"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Name is required and must be at least 3 characters long",
	}
	ErrInvalidPosition = &AppError{
		Err:            errors.New("invalid position"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Position is required and must be at least 3 characters long",
	}
	ErrInvalidHiredDate = &AppError{
		Err:            errors.New("invalid hired date"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Hired date is required and must be a valid date",
	}
	ErrInvalidSalary = &AppError{
		Err:            errors.New("invalid salary"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Salary is required and must be a valid number",
	}
	ErrEmployeeNotFound = &AppError{
		Err:            errors.New("employee not found"),
		Code:           constants.NotFoundError,
		HTTPStatusCode: http.StatusNotFound,
		PublicMsg:      "Employee not found",
	}
	ErrInvalidEmployeeId = &AppError{
		Err:            errors.New("invalid employee id"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Employee ID is required and must be a valid number",
	}
)
