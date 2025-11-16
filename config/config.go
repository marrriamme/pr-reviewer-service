package config

import (
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBConfig         *DBConfig
	ServerConfig     *ServerConfig
	MigrationsConfig *MigrationsConfig
}

func NewConfig() (*Config, error) {
	dbConfig, err := newDBConfig()
	if err != nil {
		return nil, err
	}

	serverConfig, err := newServerConfig()
	if err != nil {
		return nil, err
	}

	migrationsConfig, err := newMigrationsConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBConfig:         dbConfig,
		ServerConfig:     serverConfig,
		MigrationsConfig: migrationsConfig,
	}, nil
}

type DBConfig struct {
	User            string
	Password        string
	DB              string
	Port            int
	Host            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func newDBConfig() (*DBConfig, error) {
	user, userExists := os.LookupEnv("POSTGRES_USER")
	password, passwordExists := os.LookupEnv("POSTGRES_PASSWORD")
	dbname, dbExists := os.LookupEnv("POSTGRES_DB")
	host, hostExists := os.LookupEnv("POSTGRES_HOST")
	portStr, portExists := os.LookupEnv("POSTGRES_PORT")
	if !userExists || !passwordExists || !dbExists || !hostExists || !portExists {
		return nil, errors.New("incomplete database connection information")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid POSTGRES_PORT value")
	}

	maxOpenConns := 25
	if val, exists := os.LookupEnv("DB_MAX_OPEN_CONNS"); exists {
		if parsed, err1 := strconv.Atoi(val); err1 == nil {
			maxOpenConns = parsed
		}
	}
	maxIdleConns := 25
	if val, exists := os.LookupEnv("DB_MAX_IDLE_CONNS"); exists {
		if parsed, err2 := strconv.Atoi(val); err2 == nil {
			maxIdleConns = parsed
		}
	}
	connMaxLifetime := 5 * time.Minute
	if val, exists := os.LookupEnv("DB_CONN_MAX_LIFETIME"); exists {
		if parsed, err3 := strconv.Atoi(val); err3 == nil {
			connMaxLifetime = time.Duration(parsed) * time.Minute
		}
	}

	return &DBConfig{
		User:            user,
		Password:        password,
		DB:              dbname,
		Port:            port,
		Host:            host,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}, nil
}

func ConfigureDB(db *sql.DB, cfg *DBConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
}

type ServerConfig struct {
	Port               string
	WriteTimeout       time.Duration
	ReadTimeout        time.Duration
	IdleTimeout        time.Duration
	MaxMultipartMemory int64
}

func newServerConfig() (*ServerConfig, error) {
	port, portExist := os.LookupEnv("SERVER_PORT")
	if !portExist {
		return nil, errors.New("incomplete server configuration: missing required environment variable")
	}

	writeTimeout := getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second)
	readTimeout := getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second)
	idleTimeout := getEnvAsDuration("SERVER_IDLE_TIMEOUT", 30*time.Second)

	return &ServerConfig{
		Port:         port,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}, nil
}

type MigrationsConfig struct {
	Path string
}

func newMigrationsConfig() (*MigrationsConfig, error) {
	path, exists := os.LookupEnv("MIGRATIONS_PATH")
	if !exists {
		return nil, errors.New("MIGRATIONS_PATH is not set")
	}
	return &MigrationsConfig{
		Path: path,
	}, nil
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val, exists := os.LookupEnv(key); exists {
		if parsed, err := time.ParseDuration(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}
