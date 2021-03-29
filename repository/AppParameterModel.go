package repository

import "database/sql"

type AppParameter struct {
	Id            sql.NullInt64
	Name          sql.NullString
	Value         sql.NullString
	Description   sql.NullString
	UpdatedClient sql.NullString
	CreatedAt     sql.NullTime
	CreatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	Deleted       sql.NullBool
}
