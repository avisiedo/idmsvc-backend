package logger

import (
	"context"
	"errors"
	"time"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormLogger struct {
	slogger              *slog.Logger
	ignoreRecordNotFound bool
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.slogger.InfoContext(ctx, msg, args)
}

func (l *gormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.slogger.WarnContext(ctx, msg, args)
}

func (l *gormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.slogger.ErrorContext(ctx, msg, args)
}

func (l *gormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	elapsedTime := time.Since(begin)

	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) ||
		!l.ignoreRecordNotFound) {
		sql, rows := fc()

		l.slogger.Log(
			ctx,
			LevelTrace,
			err.Error(),
			slog.Any("error", err),
			slog.String("query", sql),
			slog.Duration("elapsed", elapsedTime),
			slog.Int64("rows", rows),
		)
	} else {
		sql, rows := fc()

		l.slogger.Log(
			ctx,
			LevelTrace,
			"SQL query executed",
			slog.String("query", sql),
			slog.Duration("elapsed", elapsedTime),
			slog.Int64("rows", rows),
		)
	}
}

func NewGormLog(ignoreRecordNotFound bool) logger.Interface {
	return &gormLogger{
		slogger:              slog.Default(),
		ignoreRecordNotFound: true,
	}
}
