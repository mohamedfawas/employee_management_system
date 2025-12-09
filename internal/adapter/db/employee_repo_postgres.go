package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
	"github.com/mohamedfawas/employee_management_system/internal/domain/repository"
	appError "github.com/mohamedfawas/employee_management_system/pkg/apperror"
)

type EmployeeRepoPostgres struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepository(pool *pgxpool.Pool) repository.EmployeeRepository {
	return &EmployeeRepoPostgres{pool: pool}
}

func (r *EmployeeRepoPostgres) CreateEmployee(ctx context.Context, employee *entity.Employee) (*entity.Employee, error) {
	query := `
		INSERT INTO employees (name, position, salary, hired_date) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, name, position, salary, hired_date, created_at
	`

	var createdEmployee entity.Employee
	row := r.pool.QueryRow(ctx, query,
		employee.Name,
		employee.Position,
		employee.Salary,
		employee.HiredDate)

	err := row.Scan(
		&createdEmployee.ID,
		&createdEmployee.Name,
		&createdEmployee.Position,
		&createdEmployee.Salary,
		&createdEmployee.HiredDate,
		&createdEmployee.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &createdEmployee, nil
}

func (r *EmployeeRepoPostgres) GetEmployeeById(ctx context.Context, id int) (*entity.Employee, error) {
	query := `
		SELECT id, name, position, salary, hired_date, created_at, updated_at 
		FROM employees 
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var employee entity.Employee
	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.Position,
		&employee.Salary,
		&employee.HiredDate,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &employee, nil
}

func (r *EmployeeRepoPostgres) GetAllEmployees(ctx context.Context) ([]*entity.Employee, error) {
	query := `
		SELECT id, name, position, salary, hired_date, created_at FROM employees
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	employees := []*entity.Employee{}
	for rows.Next() {
		var employee entity.Employee
		err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Position,
			&employee.Salary,
			&employee.HiredDate,
			&employee.CreatedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, &employee)
	}
	return employees, nil
}

func (r *EmployeeRepoPostgres) UpdateEmployee(ctx context.Context, employee *entity.Employee) (*entity.Employee, error) {
	query := `
        UPDATE employees 
        SET name = $1,
            position = $2,
            salary = $3,
            hired_date = $4,
            updated_at = NOW()
        WHERE id = $5
        RETURNING id, name, position, salary, hired_date, created_at, updated_at
    `
	row := r.pool.QueryRow(ctx, query,
		employee.Name,
		employee.Position,
		employee.Salary,
		employee.HiredDate,
		employee.ID,
	)

	var updatedEmployee entity.Employee
	err := row.Scan(
		&updatedEmployee.ID,
		&updatedEmployee.Name,
		&updatedEmployee.Position,
		&updatedEmployee.Salary,
		&updatedEmployee.HiredDate,
		&updatedEmployee.CreatedAt,
		&updatedEmployee.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &updatedEmployee, nil
}

func (r *EmployeeRepoPostgres) DeleteEmployee(ctx context.Context, id int) error {
	query := `
        DELETE FROM employees 
		WHERE id = $1
    `

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return appError.ErrEmployeeNotFound
	}
	return nil
}
