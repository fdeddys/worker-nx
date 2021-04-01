package repository

import "database/sql"

type AvailableCS struct {
	UserNexcareId sql.NullInt64
	AliasName     sql.NullString
	WorkStart     sql.NullString
	WorkEnd       sql.NullString
	BreakStart    sql.NullString
	BreakEnd      sql.NullString
	TicketId      sql.NullInt64
	Priority      sql.NullString
	Level         sql.NullString
	QueueAmount   sql.NullInt64
}