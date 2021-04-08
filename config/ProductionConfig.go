package config

import "strconv"

type ProductionConfig struct {
	Configuration
	Server struct {
		Host       string `envconfig:"MASTER_DATA_HOST"`
		Port       string `envconfig:"MASTER_DATA_PORT"`
		Version    string `json:"version"`
		ResourceID string `envconfig:"MASTER_DATA_RESOURCE_ID"`
		PrefixPath string `json:"prefix_path"`
	} `json:"server"`
	Postgresql struct {
		Address           string `envconfig:"MASTER_DATA_DB_CONNECTION"`
		Param             string `envconfig:"MASTER_DATA_DB_PARAM"`
		MaxOpenConnection int    `json:"max_open_connection"`
		MaxIdleConnection int    `json:"max_idle_connection"`
	} `json:"postgresql"`
	Email struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"server_address"`
		Secure   bool   `json:"secure"`
	} `json:"email"`
	LogFile []string `json:"log_file"`
}

func (input ProductionConfig) GetServerHost() string {
	return input.Server.Host
}
func (input ProductionConfig) GetServerPort() int {
	return convertStringParamToInt("Server Port", input.Server.Port)
}
func (input ProductionConfig) GetServerVersion() string {
	return input.Server.Version
}
func (input ProductionConfig) GetServerResourceID() string {
	return input.Server.ResourceID
}
func (input ProductionConfig) GetServerPrefixPath() string {
	return input.Server.PrefixPath
}
func (input ProductionConfig) GetPostgreSQLAddress() string {
	return input.Postgresql.Address
}
func (input ProductionConfig) GetPostgreSQLParam() string {
	return input.Postgresql.Param
}
func (input ProductionConfig) GetPostgreSQLMaxOpenConnection() int {
	return input.Postgresql.MaxOpenConnection
}
func (input ProductionConfig) GetPostgreSQLMaxIdleConnection() int {
	return input.Postgresql.MaxIdleConnection
}
func (input ProductionConfig) GetEmailUsername() string {
	return input.Email.Username
}
func (input ProductionConfig) GetEmailPassword() string {
	return input.Email.Password
}
func (input ProductionConfig) GetEmailAddress() string {
	return input.Email.Address
}
func (input ProductionConfig) GetEmailSecure() bool {
	return input.Email.Secure
}

func convertStringParamToInt(key string, value string) int {
	intPort, err := strconv.Atoi(value)
	if err != nil {
		// logModel := applicationModel.GenerateLogModel("-", "-")
		// logModel.Message = "Invalid " + key + " : " + err.Error()
		// logModel.Status = 500
		// util.LogError(logModel.ToLoggerObject())
		// os.Exit(3)
		println("Error : ", err.Error())
	}
	return intPort
}
