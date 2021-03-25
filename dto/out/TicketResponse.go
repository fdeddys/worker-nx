package out

import "time"

type Ticket struct {
	Id              int64     `json:"id"`
	Tracker         string    `json:"tracker"`
	Program         string    `json:"program"`
	ComplaintId     int64     `json:"complaint_id"`
	IssueDesc       string    `json:"issue_desc"`
	RootCauseId     int64     `json:"root_cause_id"`
	RootCauseDesc   string    `json:"root_cause_desc"`
	Status          string    `json:"status"`
	Solution        string    `json:"solution"`
	Note            string    `json:"note"`
	CustomerId      int64     `json:"customer_id"`
	LinkStorage     string    `json:"link_storage"`
	ReleaseVersion  string    `json:"release_version"`
	AssigneeId      int64     `json:"assignee_id"`
	EscalationId    int64     `json:"escalation_id"`
	RelatedTicketNo int64     `json:"related_ticket_no"`
	ContactId       int64     `json:"contact_id"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       int64     `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       int64     `json:"updated_by"`
}
