package config

type DevelopmentConfig struct {
	Configuration
	Server struct {
		Host       string `json:"host"`
		Port       int    `json:"port"`
		Version    string `json:"version"`
		ResourceID string `json:"resource_id"`
		PrefixPath string `json:"prefix_path"`
	} `json:"server"`
	Postgresql struct {
		Address           string `json:"address"`
		Param             string `json:"param"`
		MaxOpenConnection int    `json:"max_open_connection"`
		MaxIdleConnection int    `json:"max_idle_connection"`
	} `json:"postgresql"`
	Email struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"server_address"`
		Secure   bool   `json:"secure"`
	} `json:"email"`
}

func (input DevelopmentConfig) GetServerHost() string {
	return input.Server.Host
}
func (input DevelopmentConfig) GetServerPort() int {
	return input.Server.Port
}
func (input DevelopmentConfig) GetServerVersion() string {
	return input.Server.Version
}
func (input DevelopmentConfig) GetServerResourceID() string {
	return input.Server.ResourceID
}
func (input DevelopmentConfig) GetServerPrefixPath() string {
	return input.Server.PrefixPath
}
func (input DevelopmentConfig) GetPostgreSQLAddress() string {
	return input.Postgresql.Address
}
func (input DevelopmentConfig) GetPostgreSQLParam() string {
	return input.Postgresql.Param
}
func (input DevelopmentConfig) GetPostgreSQLMaxOpenConnection() int {
	return input.Postgresql.MaxOpenConnection
}
func (input DevelopmentConfig) GetPostgreSQLMaxIdleConnection() int {
	return input.Postgresql.MaxIdleConnection
}

func (input DevelopmentConfig) GetEmailUsername() string {
	return input.Email.Username
}
func (input DevelopmentConfig) GetEmailPassword() string {
	return input.Email.Password
}
func (input DevelopmentConfig) GetEmailAddress() string {
	return input.Email.Address
}
func (input DevelopmentConfig) GetEmailSecure() bool {
	return input.Email.Secure
}
