package main

import (
	"fmt"
	"nexsoft.co.id/tes-worker/config"
	"nexsoft.co.id/tes-worker/model/applicationModel"
	"nexsoft.co.id/tes-worker/serverConfig"
	"nexsoft.co.id/tes-worker/service/TaskSyncService"
	"nexsoft.co.id/tes-worker/util"
	"os"
	"os/signal"
	"strconv"

	"github.com/gobuffalo/packr/v2"
	"github.com/robfig/cron/v3"
	migrate "github.com/rubenv/sql-migrate"
	// _ "nexsoft.co.id/tes-worker/taskScheduller"
)

func main() {
	var arguments = "development"
	args := os.Args
	if len(args) > 1 {
		arguments = args[1]
	}
	fmt.Println("Argument ", arguments)
	config.GenerateConfiguration(arguments)
	serverConfig.SetServerAttribute()
	// dbMigration()

	c := cron.New()
	c.AddFunc("@every 5s", func() {
		fmt.Println("----------------> start task")
		// go taskScheduller.StartTaskScheduller()
		TaskSyncService.TaskService.DoTaskService(applicationModel.ContextModel{})
	})
	go c.Start()
	// TaskSyncService.TaskService.DoTaskService(applicationModel.ContextModel{})
	fmt.Println("Task Starting... ")
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	// ElasticSearchSyncService.ElasticSearchSyncService.DoSyncAllDBToElastic(applicationModel.ContextModel{})
	// for {

	// }

}

func dbMigration() {

	migration := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "./sql_migrations"),
	}

	if serverConfig.ServerAttribute.DBConnection != nil {
		totalRecord, errMigrate := migrate.Exec(serverConfig.ServerAttribute.DBConnection, "postgres", migration, migrate.Up)

		if errMigrate != nil {
			fmt.Println("Error db Migration !  => ", errMigrate.Error())
			os.Exit(3)
			return
		}

		logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
		logModel.Status = 200
		logModel.Message = "Applied " + strconv.Itoa(totalRecord) + " migrations!"
		util.LogInfo(logModel.ToLoggerObject())
		// util.LogInfo(logModel.ToLoggerObject())
	}
}
