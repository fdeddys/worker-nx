package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/tkanos/gonfig"
)

var ApplicationConfiguration Configuration

type Configuration interface {
	GetServerHost() string
	GetServerPort() int
	GetServerVersion() string
	GetServerResourceID() string
	GetServerPrefixPath() string
	GetPostgreSQLAddress() string
	GetPostgreSQLParam() string
	GetPostgreSQLMaxOpenConnection() int
	GetPostgreSQLMaxIdleConnection() int
	GetLogFile() []string
	GetEmailUsername() string
	GetEmailPassword() string
	GetEmailAddress() string
	GetEmailSecure() bool
}

func GenerateConfiguration(arguments string) {
	var err error
	// WorkerCoreConfig
	// enviName := os.Getenv("NexCareConfiguration")
	enviName := os.Getenv("WorkerCoreConfig")
	if arguments == "production" {
		temp := ProductionConfig{}
		err = gonfig.GetConf("config_production.json", &temp)
		if err != nil {
			fmt.Print("Errpr get config 1-> ", err)
			os.Exit(2)
		}
		err := envconfig.Process("config_development.json", &temp)
		if err != nil {
			fmt.Print("Errpr get config 2-> ", err)
			os.Exit(2)
		}
		ApplicationConfiguration = &temp
	} else {
		temp := DevelopmentConfig{}
		err = gonfig.GetConf(enviName+"/config_development.json", &temp)
		ApplicationConfiguration = &temp
	}

	if err != nil {
		fmt.Print("Errpr get config 3-> ", err)
		os.Exit(2)
	}
}
