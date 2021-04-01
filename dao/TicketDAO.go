package dao

import (
	"database/sql"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
)

type ticketDAO struct {
	FileName string
	TableName string
}

var TicketDAO = ticketDAO{}.New()

func (ticketDAO) New() (output ticketDAO) {
	output.FileName = "TicketDAO.go"
	output.TableName = "ticket"
	return
}

func (input ticketDAO) GetUnassignedTickets(db *sql.DB) (tickets []repository.Ticket, err errorModel.ErrorModel) {
	funcName := "GetUnassignedTickets"
	query := "SELECT " +
		"	ticket.id, ticket.tracker," +
		"	ticket.program_id, ticket.complaint_id," +
		"	ticket.issue_desc, ticket.root_cause_id," +
		"	ticket.root_cause_desc, ticket.status," +
		"	ticket.solution, ticket.note," +
		"	ticket.customer_id, ticket.link_storage," +
		"	ticket.release_version, ticket.assignee_id," +
		"	ticket.escalation_id, ticket.related_ticket_no," +
		"	ticket.contact_id, ticket.created_at, " +
		"	ticket.created_by, ticket.updated_at, " +
		"	ticket.updated_by, complaint_sub.cs_level, " +
		"	remark.value " +
		"FROM " + input.TableName + " " +
		"LEFT JOIN complaint_sub " +
		"	ON complaint_sub.id = ticket.complaint_id " +
		"LEFT JOIN remark " +
		"	ON complaint_sub.priority = remark.remark " +
		"WHERE " +
		"	ticket.status = 'unassign' AND " +
		"	ticket.deleted = FALSE " +
		"ORDER BY remark.value DESC, ticket.id ASC"

	rows, dbError := db.Query(query)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	for rows.Next() {
		var ticket repository.Ticket

		dbError := rows.Scan(
			&ticket.Id, &ticket.Tracker,
			&ticket.Program, &ticket.ComplaintId,
			&ticket.IssueDesc, &ticket.RootCauseId,
			&ticket.RootCauseDesc, &ticket.Status,
			&ticket.Solution, &ticket.Note,
			&ticket.CustomerId, &ticket.LinkStorage,
			&ticket.ReleaseVersion, &ticket.AssigneeId,
			&ticket.EscalationId, &ticket.RelatedTicketNo,
			&ticket.ContactId, &ticket.CreatedAt,
			&ticket.CreatedBy, &ticket.UpdatedAt,
			&ticket.UpdatedBy, &ticket.CSLevel,
			&ticket.Priority)

		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
			return
		}
		tickets = append(tickets, ticket)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ticketDAO) UpdateTicketStatus(db *sql.DB, ticket repository.Ticket) errorModel.ErrorModel {
	funcName := "UpdateTicketStatus"

	query := "UPDATE " + input.TableName + " " +
			"SET " +
			"	status = $1, " +
			"	updated_at = $2, " +
			"	updated_by = $3, " +
			"	updated_client = $4 " +
			"WHERE " +
			"	id = $5 AND " +
			"	deleted = FALSE"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	params := []interface{}{
		ticket.Status.String,
		ticket.UpdatedAt.Time,
		ticket.UpdatedBy.Int64,
		ticket.UpdatedClient.String,
		ticket.Id.Int64,
	}
	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return errorModel.GenerateNonErrorModel()
}
