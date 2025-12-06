package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
	"github.com/mohamedfawas/employee_management_system/internal/domain/repository"
)

type EmployeeRepoPostgres struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepoPostgres(pool *pgxpool.Pool) repository.EmployeeRepository {
	return &EmployeeRepoPostgres{pool: pool}
}

func (r *EmployeeRepoPostgres) CreateEmployee(ctx context.Context, employee *entity.Employee) error {
	query := `
		INSERT INTO employees (name, position, salary, hired_date) VALUES ($1, $2, $3, $4) 
		RETURNING id, name, position, salary, hired_date, created_at
	`
	row := r.pool.QueryRow(ctx, query,
		employee.Name,
		employee.Position,
		employee.Salary,
		employee.HiredDate)
	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.Position,
		&employee.Salary,
		&employee.HiredDate,
		&employee.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepoPostgres) GetEmployeeById(ctx context.Context, id int) (*entity.Employee, error) {
	query := `
		SELECT id, name, position, salary, hired_date, created_at FROM employees WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	var employee entity.Employee
	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.Position,
		&employee.Salary,
		&employee.HiredDate,
		&employee.CreatedAt)
	if err != nil {
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

func (r *EmployeeRepoPostgres) UpdateEmployee(ctx context.Context, employee *entity.Employee) error {
	query := `
        UPDATE employees 
        SET name = $1,
            position = $2,
            salary = $3,
            hired_date = $4,
            updated_at = $5
        WHERE id = $6
        RETURNING id, name, position, salary, hired_date, created_at, updated_at
    `
	row := r.pool.QueryRow(ctx, query,
		employee.Name,
		employee.Position,
		employee.Salary,
		employee.HiredDate,
		employee.UpdatedAt,
		employee.ID,
	)
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
		return err
	}
	return nil
}

func (r *EmployeeRepoPostgres) DeleteEmployee(ctx context.Context, id int) error {
	query := `
        DELETE FROM employees WHERE id = $1
    `

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("employee not found")
	}
	return nil
}
