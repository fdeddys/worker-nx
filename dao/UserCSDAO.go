package dao

import (
	"database/sql"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
)

type userCSDAO struct {
	FileName string
	TableName string
}

var UserCSDAO = userCSDAO{}.New()

func (userCSDAO) New() (output userCSDAO) {
	output.FileName = "UserCSDAO.go"
	output.TableName = "user_cs"
	return
}

func (input userCSDAO) GetAvailableCS(db *sql.DB) (listAvailableCS []repository.AvailableCS, err errorModel.ErrorModel) {
	funcName := "GetAvailableCS"

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
		" FROM " + input.TableName +
		" LEFT JOIN queue" +
		" 	ON user_cs.user_nexcare_id = queue.staff_id " +
		" INNER JOIN user_nexcare" +
		" 	ON user_cs.user_nexcare_id = user_nexcare.id AND" +
		" 	user_cs.is_locked = FALSE AND " +
		" 	user_cs.is_online = TRUE" +
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