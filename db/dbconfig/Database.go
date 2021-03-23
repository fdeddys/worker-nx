package dbconfig

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/jackc/pgx/stdlib"
)

type DBInfo struct {
	instance      *sql.DB
	driver        string
	connectionStr string
	setParams     []string
}

var instance *sql.DB
var once sync.Once

func GetDbConnection(param string, connectionString string, maxOpenConnection int, maxIdleConnection int) *sql.DB {
	_params := []string{param}
	_dbInfo := DBInfo{nil, "pgx",
		connectionString, _params}
	_db, _err := getInstance(_dbInfo)
	if _err != nil {
		fmt.Println("Error get db connection ===> ", _err)
		os.Exit(2)
	}
	_db.SetMaxOpenConns(maxOpenConnection)
	_db.SetMaxIdleConns(maxIdleConnection)
	return _db
}

func getInstance(connInfo DBInfo) (*sql.DB, error) {
	var _errOpen error
	once.Do(func() {
		dbConnStr := connInfo.connectionStr
		if connInfo.setParams != nil && len(connInfo.setParams) > 0 {

			for _, _param := range connInfo.setParams {
				dbConnStr = dbConnStr + " " + _param
			}
		}
		instance, _errOpen = sql.Open(connInfo.driver, dbConnStr)

		if _errOpen != nil {
			fmt.Printf("connect failed to DB %v", connInfo)
			instance = nil
		}
	})
	return instance, _errOpen
}
