package service

import (
	"database/sql"
	"time"

	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/model/applicationModel"
	"nexsoft.co.id/tes-worker/model/backgroundJobModel"
	"nexsoft.co.id/tes-worker/repository"
	"nexsoft.co.id/tes-worker/util"
)

func GetJobProcess(task backgroundJobModel.ChildTask, contextModel applicationModel.ContextModel, timeNow time.Time) repository.JobProcessModel {
	return repository.JobProcessModel{
		Level:         sql.NullInt32{},
		JobID:         sql.NullString{String: util.GetUUID()},
		Group:         sql.NullString{String: task.Group},
		Type:          sql.NullString{String: task.Type},
		Name:          sql.NullString{String: task.Name},
		Status:        sql.NullString{String: constant.JobProcessOnProgressStatus},
		CreatedBy:     sql.NullInt64{Int64: 0},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: "admin"},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}
}
