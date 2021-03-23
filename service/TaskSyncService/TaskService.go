package TaskSyncService

import (
	"fmt"
	"time"

	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/model/applicationModel"
	"nexsoft.co.id/tes-worker/model/backgroundJobModel"
	"nexsoft.co.id/tes-worker/repository"
	"nexsoft.co.id/tes-worker/serverConfig"
	"nexsoft.co.id/tes-worker/service"
)

type taskService struct {
	service.AbstractService
}

var TaskService = taskService{}.New()

func (input taskService) New() (output taskService) {
	output.FileName = "TaskService.go"
	return
}

func (input taskService) DoTaskService(contextModel applicationModel.ContextModel) repository.JobProcessModel {
	var listTask []backgroundJobModel.ChildTask

	fmt.Println("--> DoTaskService")
	listTask = append(listTask, input.GetSyncChildTask())
	listTask = append(listTask, input.GetSyncChildResolutionTask())

	job := service.GetJobProcess(backgroundJobModel.ChildTask{
		Group: constant.JobProcessSynchronizeGroup,
		Type:  constant.JobProcessAssignType,
		Name:  constant.JobProcessSynchronizeGroup + constant.JobProcessAssignType,
	}, contextModel, time.Now())

	job.Level.Int32 = 1
	go input.ServiceWithChildBackgroundProcess(serverConfig.ServerAttribute.DBConnection, false, listTask, job, contextModel)

	return job
}
