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
		"	ticket.updated_by " +
		"FROM " + input.TableName + " " +
		"WHERE " +
		"	ticket.status = 'unassigned' AND " +
		"	ticket.deleted = FALSE"

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
			&ticket.UpdatedBy)

		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
			return
		}
		tickets = append(tickets, ticket)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}


