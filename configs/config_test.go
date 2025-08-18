package configs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, "5432", cfg.Postgres.Port)
	assert.Equal(t, 5*time.Second, cfg.Postgres.Timeout)
	assert.Equal(t, ":8081", cfg.Server.Port)
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     AppConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid config",
			cfg:     AppConfig{Postgres: PostgresConfig{Host: "localhost", Port: "5432", User: "user", DBName: "db"}, Server: ServerConfig{Port: ":8080"}},
			wantErr: false,
		},
		{
			name:    "Empty Postgres Host",
			cfg:     AppConfig{Postgres: PostgresConfig{Host: "", Port: "5432", User: "user", DBName: "db"}, Server: ServerConfig{Port: ":8080"}},
			wantErr: true,
			errMsg:  "postgres host is required",
		},
		{
			name:    "Empty Server Port",
			cfg:     AppConfig{Postgres: PostgresConfig{Host: "localhost", Port: "5432", User: "user", DBName: "db"}, Server: ServerConfig{Port: ""}},
			wantErr: true,
			errMsg:  "server port is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	cfg, _ := Load()

	assert.Equal(t, "disable", cfg.Postgres.SSLMode)
	assert.Equal(t, 10*time.Second, cfg.Server.ReadTimeout)
}
