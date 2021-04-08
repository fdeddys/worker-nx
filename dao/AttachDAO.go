package dao

import (
	"database/sql"
	"fmt"

	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
)

type attachmentDAO struct {
	AbstractDAO
}

var AttachDAO = attachmentDAO{}.New()

func (input attachmentDAO) New() (output attachmentDAO) {
	output.FileName = "AttachmentDAO.go"
	output.TableName = "ticket_email_attachment"
	// output.ElasticSearchIndex = "bank"
	return
}

func (input attachmentDAO) InsertAttachment(tx *sql.DB, userParam repository.TicketEmailAttachment, idTicket int64) (err errorModel.ErrorModel) {
	fmt.Println("--> InsertJobProcess -->insert into table EMAIL Attachment ")

	funcName := "InsertEmailAttachment"
	query :=
		"INSERT INTO " + input.TableName + "(tiket_email_id, filename, filetype, path_cdn, updated_client, created_at, created_by, updated_at, updated_by, deleted)" +
			"	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING ID "

	param := []interface{}{
		userParam.TicketEmailID,
		userParam.Filename.String,
		userParam.Filetype.String,
		userParam.PathCdn,
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
	_, errorSt := stmt.Exec(param...)
	if errorSt != nil {
		fmt.Println("error db => ", errorS.Error())
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	// userParam.ID.Int64 = tx.
	return errorModel.GenerateNonErrorModel()
}
