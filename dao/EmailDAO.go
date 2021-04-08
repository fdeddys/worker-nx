package dao

import (
	"database/sql"
	"fmt"

	"nexsoft.co.id/tes-worker/dto/in"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
)

type emailDAO struct {
	AbstractDAO
}

var EmailDAO = emailDAO{}.New()

func (input emailDAO) New() (output emailDAO) {
	output.FileName = "EmailDAO.go"
	output.TableName = "ticket_email"
	// output.ElasticSearchIndex = "bank"
	return
}

func (input emailDAO) GetCountEmail(db *sql.DB, searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64) (result int, err errorModel.ErrorModel) {
	return 1, errorModel.ErrorModel{Code: 0}
}

func (input emailDAO) InsertEmail(tx *sql.DB, userParam repository.TicketEmail) (err errorModel.ErrorModel, ticketID int64) {
	fmt.Println("--> InsertJobProcess -->insert into table EMAIL ")

	funcName := "InsertEmail"
	query :=
		"INSERT INTO " + input.TableName + "(ticket_id, message_id, email_subject, email_date, email_sender, email_sender_name, text_plain, text_html, updated_client, created_at, created_by, updated_at, updated_by, deleted)" +
			"	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 ) RETURNING ID "

	param := []interface{}{
		userParam.TicketID.Int64,
		userParam.MessageID.String,
		userParam.EmailSubject.String,
		userParam.EmailDate.Time,
		userParam.EmailSender.String,
		userParam.EmailSenderName.String,
		userParam.TextPlain.String,
		userParam.TextHTML.String,
		userParam.UpdatedClient.String,
		userParam.CreatedAt.Time,
		userParam.CreatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdateBy.Int64,
		userParam.Deleted.Bool,
	}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		fmt.Println("error db => ", errorS.Error())
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	// var ticketID int64
	errorSt := stmt.QueryRow(param...).Scan(&ticketID)
	if errorSt != nil {
		fmt.Println("error db => ", errorS.Error())
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	// _, newId := res.()
	fmt.Println("Res insert => ", ticketID)
	// userParam.ID.Int64 = tx.
	return errorModel.GenerateNonErrorModel(), ticketID
}
