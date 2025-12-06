package entity

import "time"

type Employee struct {
	ID        int
	Name      string
	Position  string
	Salary    int
	HiredDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
