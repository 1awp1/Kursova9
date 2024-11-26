package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	PG
	Server
	Auth
	Migration
}

type PG struct {
	Username string
	Host     string
	Port     string
	DBName   string
	Password string
}

type Server struct {
	Port      string
	ReadTime  time.Duration
	WriteTime time.Duration
}

type Auth struct {
	ATDuration     time.Duration
	RFDuration     time.Duration
	PrivateKeyPath string
	PublicKeyPath  string
}

type Migration struct {
	Path string
}

func InitConfig(prefix string) (*Config, error) {
	godotenv.Load()

	conf := &Config{}
	if err := envconfig.InitWithPrefix(conf, prefix); err != nil {
		return nil, fmt.Errorf("init config error: %w", err)
	}

	return conf, nil
}
