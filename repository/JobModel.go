package repository

import "database/sql"

type JobProcessModel struct {
	ID            sql.NullInt64
	UUIDKey       sql.NullString
	ParentJobID   sql.NullString
	Level         sql.NullInt32
	JobID         sql.NullString
	Group         sql.NullString
	Type          sql.NullString
	Name          sql.NullString
	Counter       sql.NullInt32
	Total         sql.NullInt32
	Status        sql.NullString
	MessageAlert  sql.NullString
	AlertId       sql.NullString
	AlertContent  sql.NullString
	Parameter     sql.NullString
	Data          sql.NullString
	Url           sql.NullString
	Filename      sql.NullString
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	UpdatedAt     sql.NullTime
}
