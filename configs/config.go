package configs

import (
	"errors"
	"time"
)

type AppConfig struct {
	Postgres PostgresConfig
	Server   ServerConfig
}
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timeout  time.Duration
}
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (*AppConfig, error) {
	cfg := &AppConfig{
		Postgres: PostgresConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "qwerty",
			DBName:   "inventory",
			SSLMode:  "disable",
			Timeout:  5 * time.Second,
		},
		Server: ServerConfig{
			Port:         ":8081",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}
func (c *AppConfig) Validate() error {
	if c.Postgres.Host == "" {
		return errors.New("postgres host is required")
	}
	if c.Postgres.Port == "" {
		return errors.New("postgres port is required")
	}
	if c.Postgres.User == "" {
		return errors.New("postgres user is required")
	}
	if c.Postgres.DBName == "" {
		return errors.New("postgres dbname is required")
	}
	if c.Server.Port == "" {
		return errors.New("server port is required")
	}
	return nil
}
