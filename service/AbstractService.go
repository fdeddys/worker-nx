package service

import (
	"database/sql"
	"fmt"
	"time"

	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/dao"
	"nexsoft.co.id/tes-worker/model/applicationModel"
	"nexsoft.co.id/tes-worker/model/backgroundJobModel"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
	"nexsoft.co.id/tes-worker/serverConfig"
	"nexsoft.co.id/tes-worker/util"
)

var isexec bool

type AbstractService struct {
	FileName string
	Audit    bool
}

func init() {
	isexec = false
}

func (input AbstractService) ServiceWithChildBackgroundProcess(db *sql.DB, isAlertWhenError bool, listChildTask []backgroundJobModel.ChildTask, parentJob repository.JobProcessModel, contextModel applicationModel.ContextModel) {
	var err errorModel.ErrorModel

	fmt.Println("--> ServiceWithChildBackgroundProcess")

	parentJob.Total.Int32 = int32(len(listChildTask))

	err = dao.JobProcessDAO.InsertJobProcess(db, parentJob)
	if err.Error != nil {
		return
	}

	go input.DoUpdateJobEveryXMinute(constant.UpdateLastUpdateTimeInMinute, parentJob, contextModel)

	timeNow := time.Now()
	for i := 0; i < len(listChildTask); i++ {
		childJob := GetJobProcess(listChildTask[i], contextModel, timeNow)
		go input.ServiceWithBackgroundProcess(db, isAlertWhenError, parentJob, childJob, listChildTask[i], contextModel)
	}

}

func (input AbstractService) DoUpdateJobEveryXMinute(xMinute int, job repository.JobProcessModel, contextModel applicationModel.ContextModel) {

	if isexec == true {
		return
	}
	isexec = true
	fmt.Println("--> DoUpdateJobEveryXMinute-------------------------------->")

	var err errorModel.ErrorModel
	for true {
		time.Sleep(time.Duration(xMinute) * time.Minute)

		job, err = input.doUpdateJobUpdateAtOnDB(job, contextModel)
		if err.Error != nil {
			input.LogError(err, contextModel)
			isexec = false
			fmt.Println("--> DoUpdateJobEveryXMinute--------------------------------> error")
			return
		}

		if job.Status.String == constant.JobProcessDoneStatus || job.Status.String == constant.JobProcessErrorStatus {
			isexec = false
			fmt.Println("--> DoUpdateJobEveryXMinute--------------------------------> done/err")
			break
		}
	}
}

func (input AbstractService) doUpdateJobUpdateAtOnDB(job repository.JobProcessModel, contextModel applicationModel.ContextModel) (result repository.JobProcessModel, err errorModel.ErrorModel) {
	funcName := "doUpdateJobUpdateAtOnDB"
	fmt.Println("--> doUpdateJobUpdateAtOnDB  -> update status Job Proses -> selesai / error ")

	tx, errs := serverConfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	defer func() {
		if errs != nil && err.Error != nil {
			_ = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			}
			input.LogError(err, contextModel)
		} else {
			_ = tx.Commit()
		}
	}()

	result, err = dao.JobProcessDAO.GetJobProcessForUpdate(tx, job)
	if err.Error != nil {
		return
	}

	if result.Status.String == constant.JobProcessDoneStatus || result.Status.String == constant.JobProcessErrorStatus {
		return
	}

	job.UpdatedAt.Time = time.Now()

	err = dao.JobProcessDAO.UpdateJobProcessUpdateAt(tx, result)
	if err.Error != nil {
		return
	}

	return result, errorModel.GenerateNonErrorModel()
}

func (input AbstractService) LogError(err errorModel.ErrorModel, contextModel applicationModel.ContextModel) {
	contextModel.LoggerModel.Status = err.Code
	if err.CausedBy != nil {
		err.Error = err.CausedBy
	}
	// contextModel.LoggerModel.Message = util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage)
	util.LogError(contextModel.LoggerModel.ToLoggerObject())
}

func (input AbstractService) ServiceWithBackgroundProcess(db *sql.DB, isAlertWhenError bool, parentJob repository.JobProcessModel, childJob repository.JobProcessModel, task backgroundJobModel.ChildTask, contextModel applicationModel.ContextModel) {
	var err errorModel.ErrorModel
	var total int

	fmt.Println("--> ServiceWithBackgroundProcess")

	defer func() {
		if err.Error != nil {
			if isAlertWhenError {
				//todo save alert
			}

			timeNow := time.Now()
			if parentJob.JobID.String != "" {
				parentJob.Status.String = constant.JobProcessErrorStatus
				parentJob.UpdatedAt.Time = timeNow
				err = dao.JobProcessDAO.UpdateErrorJobProcess(db, parentJob)
				if err.Error != nil {
					input.LogError(err, contextModel)
				}
			}
			childJob.Status.String = constant.JobProcessErrorStatus
			childJob.UpdatedAt.Time = timeNow
			err = dao.JobProcessDAO.UpdateErrorJobProcess(db, childJob)
			if err.Error != nil {
				input.LogError(err, contextModel)
			}
		}
	}()

	if parentJob.JobID.String != "" {
		childJob.ParentJobID = parentJob.JobID
		childJob.Level.Int32 = parentJob.Level.Int32 + 1
	} else {
		childJob.Level.Int32 = 0
	}

	childJob.Counter.Int32 = 0
	total, err = task.GetCountData(db, task.Data.SearchByParam, task.Data.IsCheckStatus, task.Data.CreatedBy)
	if err.Error != nil {
		return
	}
	childJob.Total.Int32 = int32(total)
	childJob.Parameter.String = util.StructToJSON(task.Data)

	err = dao.JobProcessDAO.InsertJobProcess(db, childJob)
	if err.Error != nil {
		return
	}

	//todo kill update job process every x minute jika ditemukan error
	go input.DoUpdateJobEveryXMinute(constant.UpdateLastUpdateTimeInMinute, childJob, contextModel)

	fmt.Println("--> ServiceWithBackgroundProcess => call Task TODO")
	err = task.DoJob(db, task.Data.Data, &childJob)
	if err.Error != nil {
		return
	}

	if parentJob.JobID.String != "" {
		err = dao.JobProcessDAO.UpdateParentJobProcessCounter(db, childJob)
		if err.Error != nil {
			return
		}
	}
}
