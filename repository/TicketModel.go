package repository

import "database/sql"

type Ticket struct {
	Id              sql.NullInt64
	Tracker         sql.NullString
	Program         sql.NullString
	ComplaintId     sql.NullInt64
	CSLevel         sql.NullString
	Priority        sql.NullString
	IssueDesc       sql.NullString
	RootCauseId     sql.NullInt64
	RootCauseDesc   sql.NullString
	Status          sql.NullString
	Solution        sql.NullString
	Note            sql.NullString
	CustomerId      sql.NullInt64
	LinkStorage     sql.NullString
	ReleaseVersion  sql.NullString
	AssigneeId      sql.NullInt64
	EscalationId    sql.NullInt64
	RelatedTicketNo sql.NullInt64
	ContactId       sql.NullInt64
	TicketNo        sql.NullString
	IsReassigned    sql.NullBool
	UpdatedClient   sql.NullString
	CreatedAt       sql.NullTime
	CreatedBy       sql.NullInt64
	UpdatedAt       sql.NullTime
	UpdatedBy       sql.NullInt64
	Deleted         sql.NullBool
}
