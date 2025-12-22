package db

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SlogAdapter struct {
	Logger *slog.Logger
	Config logger.Config
}

func (a *SlogAdapter) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *a
	newLogger.Config.LogLevel = level
	return &newLogger
}

func (a *SlogAdapter) Info(ctx context.Context, msg string, data ...any) {
	if a.Config.LogLevel >= logger.Info {
		a.Logger.InfoContext(ctx, msg, "data", data)
	}
}

func (a *SlogAdapter) Warn(ctx context.Context, msg string, data ...any) {
	if a.Config.LogLevel >= logger.Warn {
		a.Logger.WarnContext(ctx, msg, "data", data)
	}
}

func (a *SlogAdapter) Error(ctx context.Context, msg string, data ...any) {
	if a.Config.LogLevel >= logger.Error {
		a.Logger.ErrorContext(ctx, msg, "data", data)
	}
}

func (a *SlogAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if a.Config.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && a.Config.LogLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		a.Logger.ErrorContext(ctx, "GORM ERROR", "err", err, "elapsed", elapsed, "rows", rows, "sql", sql)
	case elapsed > a.Config.SlowThreshold && a.Config.SlowThreshold != 0 && a.Config.LogLevel >= logger.Warn:
		a.Logger.WarnContext(ctx, "SLOW SQL", "elapsed", elapsed, "rows", rows, "sql", sql)
	case a.Config.LogLevel >= logger.Info:
		a.Logger.InfoContext(ctx, "SQL", "elapsed", elapsed, "rows", rows, "sql", sql)
	}
}
