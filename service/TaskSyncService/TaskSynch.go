package TaskSyncService

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/dao"
	"nexsoft.co.id/tes-worker/model/backgroundJobModel"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
	"nexsoft.co.id/tes-worker/serverConfig"
	"strconv"
	"time"
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
	go assignTicket(5)
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

func assignTicket(waitingTime int) {
	fmt.Println(" ----> Assign ticket")

	db := serverConfig.ServerAttribute.DBConnection
	tickets, err := dao.TicketDAO.GetUnassignedTickets(db)
	if err.Error != nil {
		print(err.Error.Error())
		return
	}

	now := time.Now()
	for _, ticket := range tickets {
		times, _ := time.Parse(time.RFC3339, now.Format("2006-01-02T15:04:05Z"))
		timeDifference := int(times.Sub(ticket.CreatedAt.Time).Minutes())

		if timeDifference > waitingTime {
			queue := repository.Queue{
				TiketId:        sql.NullInt64{Int64: ticket.Id.Int64},
				CreatedQueue:   sql.NullTime{Time: now},
				QueueStatus:    sql.NullString{String: "ready"},
				UpdatedClient:  sql.NullString{String: "3e3cb40e14d645eb8783f53a30c822d4"},
				CreatedAt:      sql.NullTime{Time: now},
				CreatedBy:      sql.NullInt64{Int64: 1},
				UpdatedAt:      sql.NullTime{Time: now},
				UpdatedBy:      sql.NullInt64{Int64: 1},
			}

			lastQueueId, err := dao.QueueDAO.InsertQueue(db, queue)
			if err.Error != nil {
				fmt.Println("-----> FAILED : Inserting ticket number " + strconv.Itoa(int(ticket.Id.Int64)))
				continue
			}

			fmt.Println("-----> SUCCESS : Inserting ticket number " + strconv.Itoa(int(ticket.Id.Int64)) + " (QueueId : "+ strconv.Itoa(int(lastQueueId)) +")")
		} else {
			fmt.Println("-----> WAITING CS TO TAKE THE TICKET FOR "+ strconv.Itoa(waitingTime) +" MINUTES...")
		}
	}
}

func reAssignTicket() {
	fmt.Println(" ----> re assign ticket")
}

func resolutionTime() {
	fmt.Println(" ----> resolution time")
}
