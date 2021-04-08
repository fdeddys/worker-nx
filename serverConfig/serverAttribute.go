package serverConfig

import (
	"crypto/tls"
	"database/sql"
	"fmt"

	pop3 "github.com/bytbox/go-pop3"
	"nexsoft.co.id/tes-worker/config"
	"nexsoft.co.id/tes-worker/db/dbconfig"
)

var ServerAttribute serverAttribute

type serverAttribute struct {
	Version      string
	DBConnection *sql.DB
	Client       *pop3.Client
}

func SetServerAttribute() {

	dbParam := config.ApplicationConfiguration.GetPostgreSQLParam()
	dbConnection := config.ApplicationConfiguration.GetPostgreSQLAddress()
	dbMaxOpenConnection := config.ApplicationConfiguration.GetPostgreSQLMaxOpenConnection()
	dbMaxIdleConnection := config.ApplicationConfiguration.GetPostgreSQLMaxIdleConnection()

	fmt.Println("SetServerAttribute ==> Param : ", dbParam, " connection ", dbConnection)
	ServerAttribute.DBConnection = dbconfig.GetDbConnection(dbParam, dbConnection, dbMaxOpenConnection, dbMaxIdleConnection)

	// ServerAttribute.Client = openEmail()
}

func openEmail() *pop3.Client {

	var client *pop3.Client
	username := config.ApplicationConfiguration.GetEmailUsername()
	password := config.ApplicationConfiguration.GetEmailPassword()
	serverAddress := config.ApplicationConfiguration.GetEmailAddress()
	secure := config.ApplicationConfiguration.GetEmailSecure()

	fmt.Println("Connect to %s\n", serverAddress)

	var dialErr error
	if secure {
		conn, err := tls.Dial("tcp", serverAddress, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			fmt.Println("error Dial ", err.Error())
		}
		client, dialErr = pop3.NewClient(conn)

	} else {
		client, dialErr = pop3.Dial(serverAddress)
	}

	if dialErr != nil {
		fmt.Println("error Dial ", dialErr.Error())
		return nil
	}

	if authErr := ServerAttribute.Client.Auth(username, password); authErr != nil {
		fmt.Println("error Auth dial ", authErr.Error())
		return nil
	}
	fmt.Println("Dial Email Success")
	return client
}
