package dao

import (
	"database/sql"
	"fmt"

	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
)

type jobProcessDAO struct {
	AbstractDAO
}

var JobProcessDAO = jobProcessDAO{}.New()

func (input jobProcessDAO) New() (output jobProcessDAO) {
	output.FileName = "JobProcessDAO.go"
	output.TableName = "job_process"
	return
}

func (input jobProcessDAO) InsertJobProcess(tx *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	fmt.Println("--> InsertJobProcess -->insert into table job_process ")
	funcName := "InsertJobProcess"
	query :=
		"INSERT INTO " + input.TableName + "(parent_job_id, level, job_id, \"group\", parameter, type, name, counter, total, created_by, created_at, created_client, updated_at) " +
			"	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) "

	param := []interface{}{
		userParam.ParentJobID.String,
		userParam.Level.Int32,
		userParam.JobID.String,
		userParam.Group.String,
		userParam.Parameter.String,
		userParam.Type.String,
		userParam.Name.String,
		userParam.Counter.Int32,
		userParam.Total.Int32,
		userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time,
		userParam.CreatedClient.String,
		userParam.UpdatedAt.Time}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		fmt.Println("error db => ", errorS.Error())
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		fmt.Println("error db => ", errorS.Error())
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) GetJobProcessForUpdate(db *sql.Tx, useParam repository.JobProcessModel) (result repository.JobProcessModel, err errorModel.ErrorModel) {
	funcName := "GetJobProcessForUpdate"
	fmt.Println("--> GetJobProcessForUpdate")

	query :=
		" SELECT " +
			"	job_id, counter, total, status " +
			" FROM " +
			input.TableName +
			" WHERE " +
			"	job_id = $1 FOR UPDATE "

	param := []interface{}{useParam.JobID.String}

	errorS := db.QueryRow(query, param...).Scan(&result.JobID, &result.Counter, &result.Total, &result.Status)

	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input jobProcessDAO) UpdateJobProcessUpdateAt(tx *sql.Tx, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateJobProcessUpdateAt"

	query := "UPDATE " + input.TableName + " SET updated_at = $1 WHERE job_id = $2 "

	param := []interface{}{userParam.UpdatedAt.Time, userParam.JobID.String}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateErrorJobProcess(tx *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateJobProcessCounter"

	query := "UPDATE " + input.TableName + " SET status = $1, updated_at = $2 WHERE job_id = $3 "

	param := []interface{}{userParam.Status.String, userParam.UpdatedAt.Time, userParam.JobID.String}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateStatusJobProcess(tx *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateStatusJobProcess"

	query := "UPDATE " + input.TableName + " SET status = $1, updated_at = $2 WHERE job_id = $3 "

	param := []interface{}{userParam.Status.String, userParam.UpdatedAt.Time, userParam.JobID.String}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateParentJobProcessCounter(db *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateParentJobProcessCounter"
	fmt.Println("--> UpdateParentJobProcessCounter")

	var parent repository.JobProcessModel
	tx, errs := db.Begin()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	defer func() {
		if errs != nil && err.Error != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	parent.JobID = userParam.ParentJobID
	parent, err = input.GetJobProcessForUpdate(tx, parent)
	if err.Error != nil {
		return
	}

	parent.Counter.Int32 += 1
	parent.UpdatedAt.Time = userParam.UpdatedAt.Time

	err = input.UpdateJobProcessCounterTx(tx, parent)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateJobProcessCounterTx(tx *sql.Tx, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateJobProcessCounterTx"
	fmt.Println("--> UpdateJobProcessCounterTx")

	if userParam.Status.String == constant.JobProcessOnProgressStatus || userParam.Status.String == constant.JobProcessOnProgressErrorStatus {
		if userParam.Counter.Int32 == userParam.Total.Int32 {
			if userParam.Status.String == constant.JobProcessOnProgressStatus {
				userParam.Status.String = constant.JobProcessDoneStatus
			} else {
				userParam.Status.String = constant.JobProcessErrorStatus
			}
		}
	}

	query := "UPDATE " + input.TableName + " SET counter = $1, status = $2, updated_at = $3 WHERE job_id = $4 "

	param := []interface{}{userParam.Counter.Int32, userParam.Status.String, userParam.UpdatedAt.Time, userParam.JobID.String}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}
