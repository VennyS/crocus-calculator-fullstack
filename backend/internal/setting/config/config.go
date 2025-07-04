package config

type Config struct {
	AppEnv string
	Server ServerConfig
}

type ServerConfig struct {
	Addr string
}
