package TaskSyncService

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/dao"
	"nexsoft.co.id/tes-worker/model/backgroundJobModel"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
	"nexsoft.co.id/tes-worker/service/AssignTicketTask"
)

// GetSyncElasticChildTask
func (input taskService) GetSyncChildTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group: constant.JobProcessSynchronizeGroup,
		Type:  constant.JobProcessAssignType,
		Name:  constant.JobProcessSyncTaskAssignTicket,
		Data: backgroundJobModel.BackgroundServiceModel{
			SearchByParam: nil,
			IsCheckStatus: false,
			CreatedBy:     0,
			Data:          nil,
		},
		GetCountData: dao.TaskDAO.GetCountTask,
		DoJob:        input.syncTask,
	}
}

// GetSyncElasticChildTask
func (input taskService) GetSyncChildResolutionTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group: constant.JobProcessSynchronizeGroup,
		Type:  constant.JobProcessResolutionTimeType,
		Name:  constant.JobProcessSyncTaskResolutionTime,
		Data: backgroundJobModel.BackgroundServiceModel{
			SearchByParam: nil,
			IsCheckStatus: false,
			CreatedBy:     0,
			Data:          nil,
		},
		GetCountData: dao.TaskDAO.GetCountTask,
		DoJob:        input.resolutionTask,
	}
}

func (input taskService) syncTask(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel) (err errorModel.ErrorModel) {
	fmt.Println(" ==================")
	fmt.Println("  TODO Assign  !!     ")
	err = AssignTicketTask.AssignTicket(db)

	// Set status Done
	childJob.Status.String = constant.JobProcessDoneStatus
	if err.Error != nil {
		childJob.Status.String = constant.JobProcessErrorStatus
	}
	err = dao.JobProcessDAO.UpdateStatusJobProcess(db, *childJob)
	fmt.Println(" ==================")
	return
}

func (input taskService) resolutionTask(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel) (err errorModel.ErrorModel) {
	fmt.Println(" ==================")
	fmt.Println("  TODO  resolution !!     ")
	go resolutionTime()
	// set status done
	fmt.Println(" ==================")
	return
}

func reAssignTicket() {
	fmt.Println(" ----> re assign ticket")
}

func resolutionTime() {
	fmt.Println(" ----> resolution time")
}
