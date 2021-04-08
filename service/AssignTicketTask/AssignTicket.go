package AssignTicketTask

import (
	"database/sql"
	"errors"
	"fmt"
	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/dao"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
	"sort"
	"strconv"
	"strings"
	"time"
)

func AssignTicket(db *sql.DB) (err errorModel.ErrorModel) {
	fmt.Println(" ----> Assign ticket")

	tickets, err := dao.TicketDAO.GetUnassignedTickets(db)
	if err.Error != nil {
		return
	}

	appParameter := repository.AppParameter{Name: sql.NullString{
		String: "tiket_waiting",
	}}
	ticketWaitingTimeParameter, err := dao.AppParameterDAO.GetParameterByName(db, appParameter)
	if err.Error != nil {
		return
	}

	ticketWaitingTime, errConv := strconv.Atoi(ticketWaitingTimeParameter.Value.String)
	if errConv != nil {
		return
	}

	listAvailableCS, err := dao.UserCSDAO.GetAvailableCS(db)
	if err.Error != nil {
		return
	}

	now := time.Now()
	for _, ticket := range tickets {
		// Check waiting time
		times, _ := time.Parse(time.RFC3339, now.Format("2006-01-02T15:04:05Z"))
		timeDifference := int(times.Sub(ticket.CreatedAt.Time).Minutes())

		if timeDifference >= ticketWaitingTime {
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

			// Finding Available CS
			csIndex, availableCs, errs := findAvailableCS(listAvailableCS, ticket, now)
			if errs != nil {
				print(errs.Error())
				continue
			}
			queue.StaffId.Int64 = availableCs.UserNexcareId.Int64

			// Is there any CS available ?
			if queue.StaffId.Int64 > 0 && csIndex > -1 {
				listAvailableCS[csIndex].QueueAmount.Int64 += 1
				listAvailableCS[csIndex].TicketId.Int64 = ticket.Id.Int64

				// Insert to table queue
				err = insertTicketToQueue(db, queue)
				if err.Error != nil {
					continue
				}

				// Update ticket status to 'assigned'
				err = updateTicketStatus(db, ticket, now)
				if err.Error != nil {
					continue
				}

				// Re-sort List CS
				reSortListCS(listAvailableCS)
			} else {
				fmt.Println("-----> NOT FOUND : Cannot find available CS ")
			}
		} else {
			fmt.Println("-----> WAITING CS TO TAKE THE TICKET FOR " + strconv.Itoa(ticketWaitingTime) + " MINUTES...")
		}
	}
	return
}

func findAvailableCS(listCS []repository.AvailableCS, ticket repository.Ticket, now time.Time) (index int, availableCS repository.AvailableCS, err error) {
	for i, cs := range listCS {
		// Check CS Level
		if cs.Level.String != ticket.CSLevel.String {
			index = -1
			continue
		}

		// Check Shift
		isAvailable, err := isCSAvailableInShift(cs, now)
		if err != nil {
			index = -1
			fmt.Println("-----> FAILED : " + err.Error())
			continue
		}
		if !isAvailable {
			index = -1
			continue
		}
		index = i
		availableCS = cs
		break
	}
	return
}

func insertTicketToQueue(db *sql.DB, queue repository.Queue) (err errorModel.ErrorModel) {
	lastQueueId, err := dao.QueueDAO.InsertQueue(db, queue)
	if err.Error != nil {
		fmt.Println("-----> FAILED : Inserting ticket number " + strconv.Itoa(int(queue.TiketId.Int64)))
		fmt.Println(err.CausedBy.Error())
		return
	}
	fmt.Println("-----> SUCCESS : Inserting ticket number " + strconv.Itoa(int(queue.TiketId.Int64)) + " (QueueId : " + strconv.Itoa(int(lastQueueId)) + ")")
	return
}

func updateTicketStatus(db *sql.DB, ticket repository.Ticket, now time.Time) (err errorModel.ErrorModel) {
	err = dao.TicketDAO.UpdateTicketStatus(db, repository.Ticket{
		Id: ticket.Id,
		Status: sql.NullString{String: "assigned"},
		UpdatedAt: sql.NullTime{Time: now},
		UpdatedBy: sql.NullInt64{Int64: 1},
		UpdatedClient: sql.NullString{String: "3e3cb40e14d645eb8783f53a30c822d4"},
	})

	if err.Error != nil {
		fmt.Println("--------> FAILED : Updating ticket status failed")
		fmt.Println(err.CausedBy.Error())
		return
	}

	fmt.Println("--------> SUCCESS : Successfully updating ticket status")
	return
}

func reSortListCS(listAvailableCS []repository.AvailableCS) {
	sort.SliceStable(listAvailableCS, func(i, j int) bool {
		return listAvailableCS[i].QueueAmount.Int64 < listAvailableCS[j].QueueAmount.Int64
	})
	sort.SliceStable(listAvailableCS, func(i, j int) bool {
		return listAvailableCS[i].TicketId.Int64 < listAvailableCS[j].TicketId.Int64
	})
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