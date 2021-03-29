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

func (input queueDAO) GetAvailableCS(db *sql.DB) (listAvailableCS []repository.AvailableCS, err errorModel.ErrorModel) {
	funcName := "GetAvailableCS"

	//query := "SELECT " +
	//		"user_cs.user_nexcare_id," +
	//		"user_nexcare.alias_name," +
	//		"queue.tiket_id," +
	//		"remark.value," +
	//		"user_cs.cs_level," +
	//		"count(tiket_id) AS queue_amount " +
	//	"FROM user_cs " +
	//	"INNER JOIN user_nexcare " +
	//	"	ON user_cs.user_nexcare_id = user_nexcare.id AND " +
	//	"	   user_cs.is_locked = FALSE AND " +
	//	"	   user_cs.is_online = TRUE " +
	//	"LEFT JOIN queue " +
	//	"	ON user_cs.user_nexcare_id = queue.staff_id " +
	//	"LEFT JOIN ticket " +
	//	"	ON queue.tiket_id = ticket.id " +
	//	"LEFT JOIN complaint_sub " +
	//	"	ON complaint_sub.id = ticket.complaint_id " +
	//	"LEFT JOIN remark " +
	//	"	ON complaint_sub.priority = remark.remark " +
	//	"GROUP BY " +
	//	"	user_cs.user_nexcare_id," +
	//	"	user_nexcare.alias_name," +
	//	"	queue.tiket_id,"+
	//	"	remark.value,"+
	//	"	user_cs.cs_level " +
	//	"ORDER BY queue_amount ASC, user_nexcare.alias_name ASC"

	query := "SELECT " +
		" 	user_cs.user_nexcare_id," +
		" 	user_nexcare.alias_name," +
		"	list_shift.work_start," +
		"	list_shift.work_end," +
		"	list_shift.break_start," +
		"	list_shift.break_end," +
		" 	user_cs.cs_level," +
		" 	MAX(queue.tiket_id) ticket_id," +
		" 	COUNT(staff_id) AS queue_amount" +
		" FROM user_cs" +
		" LEFT JOIN queue" +
		" 	ON user_cs.user_nexcare_id = queue.staff_id AND" +
		" 	user_cs.is_locked = FALSE AND " +
		" 	user_cs.is_online = TRUE" +
		" INNER JOIN user_nexcare" +
		" 	ON user_cs.user_nexcare_id = user_nexcare.id" +
		" LEFT JOIN ticket" +
		" 	ON queue.tiket_id = ticket.id" +
		" LEFT JOIN list_shift" +
		"	ON user_cs.shift_id = list_shift.id" +
		" GROUP BY " +
		" 	staff_id, " +
		" 	user_cs.user_nexcare_id," +
		" 	user_nexcare.alias_name," +
		" 	user_cs.cs_level," +
		"	list_shift.work_start," +
		"	list_shift.work_end," +
		"	list_shift.break_start," +
		"	list_shift.break_end " +
		" ORDER BY queue_amount ASC, ticket_id ASC, alias_name ASC"
	rows, dbError := db.Query(query)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	for rows.Next() {
		var availableCS repository.AvailableCS
		dbError := rows.Scan(
			&availableCS.UserNexcareId,
			&availableCS.AliasName,
			&availableCS.WorkStart,
			&availableCS.WorkEnd,
			&availableCS.BreakStart,
			&availableCS.BreakEnd,
			&availableCS.Level,
			&availableCS.TicketId,
			&availableCS.QueueAmount)
		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
			return
		}
		listAvailableCS = append(listAvailableCS, availableCS)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}