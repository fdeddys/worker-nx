package serverConfig

import (
	"database/sql"
	"fmt"

	"nexsoft.co.id/tes-worker/config"
	"nexsoft.co.id/tes-worker/db/dbconfig"
)

var ServerAttribute serverAttribute

type serverAttribute struct {
	Version      string
	DBConnection *sql.DB
}

func SetServerAttribute() {

	dbParam := config.ApplicationConfiguration.GetPostgreSQLParam()
	dbConnection := config.ApplicationConfiguration.GetPostgreSQLAddress()
	dbMaxOpenConnection := config.ApplicationConfiguration.GetPostgreSQLMaxOpenConnection()
	dbMaxIdleConnection := config.ApplicationConfiguration.GetPostgreSQLMaxIdleConnection()

	fmt.Println("SetServerAttribute ==> Param : ", dbParam, " connection ", dbConnection)
	ServerAttribute.DBConnection = dbconfig.GetDbConnection(dbParam, dbConnection, dbMaxOpenConnection, dbMaxIdleConnection)

}
