package configs

import (
	"errors"
	"log"
	"os"
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

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func Load() (*AppConfig, error) {
	cfg := &AppConfig{
		Postgres: PostgresConfig{
			Host:     os.Getenv("DB_HOST"), // Только из переменных окружения
			Port:     os.Getenv("DB_PORT"), // без fallback значений!
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timeout:  5 * time.Second,
		},
		Server: ServerConfig{
			Port:         ":" + getEnv("SERVER_PORT", "8081"),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	log.Printf("Final config: %+v", cfg)
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
