package constant

const JobProcessDoneStatus = "OK"
const JobProcessErrorStatus = "ERROR"
const JobProcessOnProgressStatus = "ONPROGRESS"
const JobProcessOnProgressErrorStatus = "ONPROGRESS-ERROR"
const JobProcessSynchronizeGroup = "Synchronize"
const JobProcessAssignType = "Assign Ticket"
const JobProcessResolutionTimeType = "ResolutionTIme"

const JobProcessSyncTaskAssignTicket = JobProcessSynchronizeGroup + " " + JobProcessAssignType + " Task scheduller"
const JobProcessSyncTaskResolutionTime = JobProcessSynchronizeGroup + " " + JobProcessResolutionTimeType + " Task scheduller"

const UpdateLastUpdateTimeInMinute = 1
