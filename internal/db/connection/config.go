package connection

type DatabaseConfig struct {
	Host     string
	Username string
	Password string
	DBName   string
	Port     string
	AppName  string
	SSLMode  string
	Timezone string
}
