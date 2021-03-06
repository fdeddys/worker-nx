package TaskSyncService

import (
	"database/sql"
	"fmt"

	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/dao"
	"nexsoft.co.id/tes-worker/model/backgroundJobModel"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
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
	go assignTicket()
	// Set status Done
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

func assignTicket() {
	fmt.Println(" ----> Assign ticket")
}

func reAssignTicket() {
	fmt.Println(" ----> re assign ticket")
}

func resolutionTime() {
	fmt.Println(" ----> resolution time")
}
