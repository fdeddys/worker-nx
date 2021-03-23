package dao

import (
	"database/sql"

	"nexsoft.co.id/tes-worker/dto/in"
	"nexsoft.co.id/tes-worker/model/errorModel"
)

type taskDAO struct {
	AbstractDAO
}

var TaskDAO = taskDAO{}.New()

func (input taskDAO) New() (output taskDAO) {
	output.FileName = "TaskDAO.go"
	output.TableName = "task"
	// output.ElasticSearchIndex = "bank"
	return
}

func (input taskDAO) GetCountTask(db *sql.DB, searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64) (result int, err errorModel.ErrorModel) {
	return 1, errorModel.ErrorModel{Code: 0}
}
