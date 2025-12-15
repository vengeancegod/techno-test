package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	loggerLevelEnvName      = "LOGGER_LEVEL"
	loggerNoColorEnvName    = "LOGGER_NO_COLOR"
	loggerTimeFormatEnvName = "LOGGER_TIME_FORMAT"
	loggerTimeLocationEnv   = "LOGGER_TIME_LOCATION"
	loggerLogsDirEnvName    = "LOGGER_LOGS_DIR"
	loggerToStdoutEnvName   = "LOGGER_TO_STDOUT"
	loggerFileNameEnvName   = "LOGGER_FILE_NAME"
)

type LoggerConfig interface {
	Level() string
	NoColor() bool
	TimeFormat() string
	TimeLocation() string
	EnableFileLog() bool
	LogsDir() string
	LogFileName() string
	Initialize() error
}

type loggerConfig struct {
	level         string
	noColor       bool
	timeFormat    string
	timeLocation  string
	enableFileLog bool
	logsDir       string
	logFileName   string
	logToFile     bool
	logToStdout   bool
}

func NewLoggerConfig() (LoggerConfig, error) {
	level := getEnvOrDefault(loggerLevelEnvName, "info")
	noColor := getEnvOrDefault(loggerNoColorEnvName, "false") == "true"
	timeFormat := getEnvOrDefault(loggerTimeFormatEnvName, "2006-01-02 15:04:05")
	timeLocation := getEnvOrDefault(loggerTimeLocationEnv, "Local")
	logsDir := getEnvOrDefault(loggerLogsDirEnvName, "./logs")
	logFileName := getEnvOrDefault(loggerFileNameEnvName, "taskmanager.log")

	cfg := &loggerConfig{
		level:        level,
		noColor:      noColor,
		timeFormat:   timeFormat,
		timeLocation: timeLocation,
		logsDir:      logsDir,
		logFileName:  logFileName,
	}

	return cfg, nil
}

func (c *loggerConfig) Level() string {
	return c.level
}

func (c *loggerConfig) NoColor() bool {
	return c.noColor
}

func (c *loggerConfig) TimeFormat() string {
	return c.timeFormat
}

func (c *loggerConfig) TimeLocation() string {
	return c.timeLocation
}

func (c *loggerConfig) EnableFileLog() bool {
	return c.enableFileLog
}

func (c *loggerConfig) LogsDir() string {
	return c.logsDir
}

func (c *loggerConfig) LogFileName() string {
	return c.logFileName
}

func (c *loggerConfig) Initialize() error {
	level, err := zerolog.ParseLevel(strings.ToLower(c.level))
	if err != nil {
		return fmt.Errorf("invalid log level '%s': %w", c.level, err)
	}
	zerolog.SetGlobalLevel(level)

	loc, err := time.LoadLocation(c.timeLocation)
	if err != nil {
		loc = time.Local
	}
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(loc)
	}

	zerolog.TimeFieldFormat = c.timeFormat

	var writers []io.Writer

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: c.timeFormat,
		NoColor:    c.noColor,
	}
	writers = append(writers, consoleWriter)

	logFilePath := filepath.Join(c.logsDir, c.logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	writers = append(writers, logFile)

	multiWriter := io.MultiWriter(writers...)
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	return nil
}

func GetLogger(component string) zerolog.Logger {
	return log.Logger.With().Str("component", component).Logger()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
