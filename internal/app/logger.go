package app

import (
	"log/slog"
	"os"
)

func SetupLogger(level string) *slog.Logger {
	logLevel, err := parseLevel(level)
	if err != nil {
		logLevel = slog.LevelDebug
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
	)
	slog.SetDefault(log)

	return log
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
