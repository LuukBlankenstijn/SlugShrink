//go:build prod

package logging

import (
	"log/slog"
	"os"
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     readLevelFromEnv(),
	}))
}
