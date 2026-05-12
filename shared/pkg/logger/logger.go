package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func Get() *zerolog.Logger {
	return &log.Logger
}

func Info(msg string, fields ...interface{}) {
	log.Info().Fields(fields).Msg(msg)
}

func Error(msg string, err error, fields ...interface{}) {
	log.Error().Err(err).Fields(fields).Msg(msg)
}

func Debug(msg string, fields ...interface{}) {
	log.Debug().Fields(fields).Msg(msg)
}

func Warn(msg string, fields ...interface{}) {
	log.Warn().Fields(fields).Msg(msg)
}
