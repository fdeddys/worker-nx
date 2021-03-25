package repository

import "database/sql"

type Queue struct {
	Id             sql.NullInt64
	StaffId        sql.NullInt64
	TiketId        sql.NullInt64
	LevelId        sql.NullInt64
	CreatedQueue   sql.NullTime
	QueueStatus    sql.NullString
	StartExec      sql.NullTime
	DoneExec       sql.NullTime
	IsOpen         sql.NullBool
	ResolutionTime sql.NullInt64
	ResponseTime   sql.NullInt64
	ResponseTimeBy sql.NullInt64
	UpdatedClient  sql.NullString
	CreatedAt      sql.NullTime
	CreatedBy      sql.NullInt64
	UpdatedAt      sql.NullTime
	UpdatedBy      sql.NullInt64
	Deleted        sql.NullBool
}