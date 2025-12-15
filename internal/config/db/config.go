package db

import (
	"errors"
	"os"
)

const (
	dbHostEnvName = "PSQL_HOST"
	dbPortEnvName = "PSQL_PORT"
	dbUserEnvName = "PSQL_USER"
	dbPassEnvName = "PSQL_PASSWORD"
	dbNameEnvName = "PSQL_DB"
	dbSSLEnvName  = "PSQL_SSLMODE"
)

type DBConfig interface {
	Host() string
	Port() string
	User() string
	Password() string
	DBName() string
	SSLMode() string
	ConnectionString() string
}

type dbConfig struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
	sslMode  string
}

func NewDBConfig() (DBConfig, error) {
	host := os.Getenv(dbHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("database host not found")
	}

	port := os.Getenv(dbPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("database port not found")
	}

	user := os.Getenv(dbUserEnvName)
	if len(user) == 0 {
		return nil, errors.New("database user not found")
	}

	password := os.Getenv(dbPassEnvName)
	if len(password) == 0 {
		return nil, errors.New("database password not found")
	}

	dbName := os.Getenv(dbNameEnvName)
	if len(dbName) == 0 {
		return nil, errors.New("database name not found")
	}

	sslMode := os.Getenv(dbSSLEnvName)
	if len(sslMode) == 0 {
		sslMode = "disable"
	}

	return &dbConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbName:   dbName,
		sslMode:  sslMode,
	}, nil
}

func (c *dbConfig) Host() string {
	return c.host
}

func (c *dbConfig) Port() string {
	return c.port
}

func (c *dbConfig) User() string {
	return c.user
}

func (c *dbConfig) Password() string {
	return c.password
}

func (c *dbConfig) DBName() string {
	return c.dbName
}

func (c *dbConfig) SSLMode() string {
	return c.sslMode
}

func (c *dbConfig) ConnectionString() string {
	return "host=" + c.host +
		" port=" + c.port +
		" user=" + c.user +
		" password=" + c.password +
		" dbname=" + c.dbName +
		" sslmode=" + c.sslMode
}
