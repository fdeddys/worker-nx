package dao

import (
	"database/sql"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
)

type appParameterDAO struct {
	FileName string
	TableName string
}

var AppParameterDAO = appParameterDAO{}.New()

func (appParameterDAO) New() (output appParameterDAO) {
	output.FileName = "AppParameterDAO.go"
	output.TableName = "app_parameter"
	return
}

func (input appParameterDAO) GetParameterByName(db *sql.DB, appParameter repository.AppParameter) (result repository.AppParameter, err errorModel.ErrorModel) {
	funcName := "GetParameterByName"

	query := "SELECT " +
			"	id, name, value, description " +
			"FROM " + input.TableName + " " +
			"WHERE " +
			"	name = $1 AND " +
			"	deleted = FALSE"

	results := db.QueryRow(query, appParameter.Name.String)
	dbError := results.Scan(&result.Id.Int64, &result.Name.String, &result.Value.String, &result.Description.String)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}