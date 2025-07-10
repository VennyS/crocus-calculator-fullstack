package config

type Config struct {
	AppEnv string
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Addr string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}
