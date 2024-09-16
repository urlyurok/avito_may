package entity

import "time"

type Organizaton struct {
	Id          int       `db:"id"`
	Name        string    `db:"balance"`
	Description string    `db:"description"`
	Type        string    `db:"service_type"`
	CreatedAt   time.Time `db:"status"`
	UpdatedAt   time.Time `db:"organization_id"`
}
