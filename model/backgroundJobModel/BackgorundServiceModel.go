package backgroundJobModel

import (
	"database/sql"

	"nexsoft.co.id/tes-worker/dto/in"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
)

type ChildTask struct {
	Group        string
	Type         string
	Name         string
	Data         BackgroundServiceModel
	GetCountData func(*sql.DB, []in.SearchByParam, bool, int64) (int, errorModel.ErrorModel)
	DoJob        func(*sql.DB, interface{}, *repository.JobProcessModel) errorModel.ErrorModel
}

type BackgroundServiceModel struct {
	SearchByParam []in.SearchByParam
	IsCheckStatus bool
	CreatedBy     int64
	Data          interface{}
}
