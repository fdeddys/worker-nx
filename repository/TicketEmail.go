package repository

import "database/sql"

type TicketEmail struct {
	ID              sql.NullInt64
	TicketID        sql.NullInt64
	MessageID       sql.NullString
	EmailSubject    sql.NullString
	EmailDate       sql.NullTime
	EmailSender     sql.NullString
	EmailSenderName sql.NullString
	TextPlain       sql.NullString
	TextHTML        sql.NullString
	UpdatedClient   sql.NullString
	CreatedBy       sql.NullInt64
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
	UpdateBy        sql.NullInt64
	Deleted         sql.NullBool
}

type TicketEmailAttachment struct {
	ID            sql.NullInt64
	TicketEmailID sql.NullInt64
	Filename      sql.NullString
	Filetype      sql.NullString
	PathCdn       sql.NullString
	UpdatedClient sql.NullString
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	UpdatedAt     sql.NullTime
	UpdateBy      sql.NullInt64
	Deleted       sql.NullBool
}
