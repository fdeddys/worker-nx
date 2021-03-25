package dao

import (
	"database/sql"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
)

type queueDAO struct {
	AbstractDAO
}

var QueueDAO = queueDAO{}.New()

func (input queueDAO) New() (output queueDAO) {
	output.FileName = "Queue.go"
	output.TableName = "queue"
	return
}

func (input queueDAO) InsertQueue(db *sql.DB, queue repository.Queue) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertQueue"

	query := "INSERT INTO " + input.TableName + "(" +
		"	staff_id," +
		"	tiket_id," +
		"	level_id," +
		"	created_queue," +
		"	queue_status," +
		"	start_exec," +
		"	done_exec," +
		"	is_open," +
		"	resolution_time," +
		"	response_time," +
		"	response_time_by," +
		"	updated_client," +
		"	created_at," +
		"	created_by," +
		"	updated_at," +
		"	updated_by" +
		") " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) " +
		"RETURNING id"

	params := []interface{}{
		queue.StaffId.Int64,
		queue.TiketId.Int64,
		queue.LevelId.Int64,
		queue.CreatedQueue.Time,
		queue.QueueStatus.String,
		queue.StartExec.Time,
		queue.DoneExec.Time,
		queue.IsOpen.Bool,
		queue.ResolutionTime.Int64,
		queue.ResponseTime.Int64,
		queue.ResponseTimeBy.Int64,
		queue.UpdatedClient.String,
		queue.CreatedAt.Time,
		queue.CreatedBy.Int64,
		queue.UpdatedAt.Time,
		queue.UpdatedBy.Int64}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}