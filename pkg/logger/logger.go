package logger

import (
	"log/slog"
	"os"
)

// TODO find best practice for logger and make it better in this application
// TODO pass the file path to New function instead of hardcoding it in code .
func New() *slog.Logger {
	// TODO: read this article : https://gist.github.com/bmcculley/055518d703899ec3a91f7aae732b04d2
	logFile, err := os.OpenFile("/home/amir/Desktop/ShareExpense/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("failed to open log file")
	}
	loggerHandler := slog.NewJSONHandler(logFile,
		&slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	loggerInstance := slog.New(loggerHandler)
	return loggerInstance

}
