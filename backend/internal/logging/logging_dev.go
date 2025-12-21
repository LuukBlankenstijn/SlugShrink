//go:build !prod

package logging

import (
	"log/slog"
	"os"
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     readLevelFromEnv(),
	}))
}
