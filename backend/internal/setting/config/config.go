package config

import (
	"fmt"
	"time"
)

type Config struct {
	AppEnv string
	Server ServerConfig
	DB     DBConfig
	HTTP   HTTPConfig
}

type HTTPConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type ServerConfig struct {
	Addr string
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		db.Driver,
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)
}
