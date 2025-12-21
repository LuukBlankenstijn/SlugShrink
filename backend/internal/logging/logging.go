package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/LuukBlankenstijn/gewish/internal/utils"
)

var logger = newLogger()

func init() {
	slog.SetDefault(logger)
}

func Logger() *slog.Logger {
	return logger
}

func Fatal(msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}

func readLevelFromEnv() slog.Level {
	value := strings.ToLower(utils.EnvOrDefault("LOG_LEVEL", "info"))
	switch value {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
