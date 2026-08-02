package logger

import (
	"log/slog"
	"os"
	"sync"
)

type MyLogger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

var (
	loggerInstance *slog.Logger
	once           sync.Once
)

func NewLogger() *slog.Logger {
	once.Do(func() {
		// TODO: read this article : https://gist.github.com/bmcculley/055518d703899ec3a91f7aae732b04d2
		logFile, err := os.OpenFile("/home/amir/Desktop/ShareExpense/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("failed to open log file")
		}
		loggerHandler := slog.NewJSONHandler(logFile,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
		loggerInstance = slog.New(loggerHandler)
	})
	return loggerInstance
}
