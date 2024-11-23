package config

import (
	"fmt"
	"time"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	PG
	Server
	Auth
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

func InitConfig(prefix string) (*Config, error) {
	conf := &Config{}
	if err := envconfig.InitWithPrefix(conf, prefix); err != nil {
		return nil, fmt.Errorf("init config error: %w", err)
	}

	return conf, nil
}
