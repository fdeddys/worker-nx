package TaskSyncService

import (
	"database/sql"
	"errors"
	"fmt"
	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/dao"
	"nexsoft.co.id/tes-worker/model/backgroundJobModel"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
	"sort"
	"strconv"
	"strings"
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
	go assignTicket(db)
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

func assignTicket(db *sql.DB) {
	fmt.Println(" ----> Assign ticket")

	tickets, err := dao.TicketDAO.GetUnassignedTickets(db)
	if err.Error != nil {
		print(err.Error.Error())
		return
	}

	appParameter := repository.AppParameter{Name: sql.NullString{
		String: "tiket_waiting",
	}}
	ticketWaitingTimeParameter, err := dao.AppParameterDAO.GetParameterByName(db, appParameter)
	if err.Error != nil {
		print(err.CausedBy.Error())
		return
	}

	ticketWaitingTime, errConv := strconv.Atoi(ticketWaitingTimeParameter.Value.String)
	if errConv != nil {
		print(err.CausedBy.Error())
		return
	}

	listAvailableCS, err := dao.QueueDAO.GetAvailableCS(db)
	if err.Error != nil {
		print(err.CausedBy.Error())
		return
	}

	now := time.Now()
	for _, ticket := range tickets {
		times, _ := time.Parse(time.RFC3339, now.Format("2006-01-02T15:04:05Z"))
		timeDifference := int(times.Sub(ticket.CreatedAt.Time).Minutes())

		if timeDifference > ticketWaitingTime {
			queue := repository.Queue{
				TiketId:       sql.NullInt64{Int64: ticket.Id.Int64},
				CreatedQueue:  sql.NullTime{Time: now},
				QueueStatus:   sql.NullString{String: "ready"},
				UpdatedClient: sql.NullString{String: "3e3cb40e14d645eb8783f53a30c822d4"},
				CreatedAt:     sql.NullTime{Time: now},
				CreatedBy:     sql.NullInt64{Int64: 1},
				UpdatedAt:     sql.NullTime{Time: now},
				UpdatedBy:     sql.NullInt64{Int64: 1},
			}

			for i, cs := range listAvailableCS {
				if cs.Level.String != ticket.CSLevel.String {
					continue
				}

				isAvailable, errs := isCSAvailableInShift(cs, now)
				if errs != nil {
					fmt.Println("-----> FAILED : " + errs.Error())
					continue
				}
				if !isAvailable {
					continue
				}

				listAvailableCS[i].QueueAmount.Int64 += 1
				listAvailableCS[i].TicketId.Int64 = ticket.Id.Int64
				sort.SliceStable(listAvailableCS, func(i, j int) bool {
					return listAvailableCS[i].QueueAmount.Int64 < listAvailableCS[j].QueueAmount.Int64
				})
				sort.SliceStable(listAvailableCS, func(i, j int) bool {
					return listAvailableCS[i].TicketId.Int64 < listAvailableCS[j].TicketId.Int64
				})
				queue.StaffId.Int64 = cs.UserNexcareId.Int64
				break
			}

			if queue.StaffId.Int64 > 0 {
				lastQueueId, err := dao.QueueDAO.InsertQueue(db, queue)
				if err.Error != nil {
					fmt.Println("-----> FAILED : Inserting ticket number " + strconv.Itoa(int(ticket.Id.Int64)))
					continue
				}

				fmt.Println("-----> SUCCESS : Inserting ticket number " + strconv.Itoa(int(ticket.Id.Int64)) + " (QueueId : " + strconv.Itoa(int(lastQueueId)) + ")")
			} else {
				fmt.Println("-----> NOT FOUND : Cannot find available CS ")
			}
		} else {
			fmt.Println("-----> WAITING CS TO TAKE THE TICKET FOR " + strconv.Itoa(ticketWaitingTime) + " MINUTES...")
		}
	}
}

func isCSAvailableInShift(cs repository.AvailableCS, now time.Time) (result bool, err error) {
	workStart, err := generateClock(cs.WorkStart.String, now)
	if err != nil {
		return
	}

	workEnd, err := generateClock(cs.WorkEnd.String, now)
	if err != nil {
		return
	}

	breakStart, err := generateClock(cs.BreakStart.String, now)
	if err != nil {
		return
	}

	breakEnd, err := generateClock(cs.BreakEnd.String, now)
	if err != nil {
		return
	}

	isBreak := now.After(breakStart) && now.Before(breakEnd)
	if isBreak {
		return false, nil
	}

	return now.After(workStart) && now.Before(workEnd), nil
}

func generateClock(strClock string, now time.Time) (result time.Time, err error) {
	clock := strings.Split(strClock, constant.HourSeparator)
	if len(clock) == 2 {
		hour, errHour := strconv.Atoi(clock[0])
		minute, errMinute := strconv.Atoi(clock[1])

		if errHour != nil {
			err = errHour
			return
		}

		if errMinute != nil {
			err = errMinute
			return
		}

		return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location()), nil
	}
	return time.Time{}, errors.New("clock format is invalid")
}

func assignToAvailableCS(listAvailableCS []repository.AvailableCS, ticket repository.Ticket, now time.Time) (index int, queue repository.Queue, isFound bool) {
	queue = repository.Queue{
		TiketId:       sql.NullInt64{Int64: ticket.Id.Int64},
		CreatedQueue:  sql.NullTime{Time: now},
		QueueStatus:   sql.NullString{String: "ready"},
		UpdatedClient: sql.NullString{String: "3e3cb40e14d645eb8783f53a30c822d4"},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedBy:     sql.NullInt64{Int64: 1},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: 1},
	}

	j := -1
	for i, availableCS := range listAvailableCS {
		if ticket.CSLevel.String != availableCS.Level.String {
			continue
		}

		queue.StaffId.Int64 = availableCS.UserNexcareId.Int64
		isFound = true
		j = i
		break
	}
	index = j
	return
}

//func findMinimumQueueAmount(listAvailableCS []repository.AvailableCS) (result repository.AvailableCS, isFound bool) {
//	if len(listAvailableCS) == 0 {
//		return
//	}
//
//	availableCS := repository.AvailableCS{}
//	for _, item := range listAvailableCS {
//		if item.QueueAmount.Int64 < listAvailableCS[0].QueueAmount.Int64 {
//			availableCS = item
//		}
//	}
//
//	if availableCS.UserNexcareId.Int64 < 1 {
//		return
//	}
//
//	result = availableCS
//	return availableCS, true
//}

func reAssignTicket() {
	fmt.Println(" ----> re assign ticket")
}

func resolutionTime() {
	fmt.Println(" ----> resolution time")
}
