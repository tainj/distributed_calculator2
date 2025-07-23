package logger

import (
	"context"
	"log/slog"
)

const (
	LoggerKey   = "logger"
	RequestID   = "requestID"
	ServiceName = "service"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...any)
	Error(ctx context.Context, msg string, fields ...any)
}

type logger struct {
	serviceName string
	logger      *slog.Logger
}

func (l *logger) Info(ctx context.Context, msg string, fields ...any) {
	// добавляем serviceName и requestID в логи
	attrs := []any{slog.String(ServiceName, l.serviceName)}
	if ctx.Value(RequestID) != nil {
		attrs = append(attrs, slog.String(RequestID, ctx.Value(RequestID).(string)))
	}

	// логируем с контекстом
	l.logger.InfoContext(ctx, msg, append(attrs, fields...)...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...any) {
	// добавляем serviceName и requestID в логи
	attrs := []any{slog.String(ServiceName, l.serviceName)}
	if ctx.Value(RequestID) != nil {
		attrs = append(attrs, slog.String(RequestID, ctx.Value(RequestID).(string)))
	}

	// логируем с контекстом
	l.logger.ErrorContext(ctx, msg, append(attrs, fields...)...)
}

// New создаёт новый логгер
func New(serviceName string) Logger {
	return &logger{
		serviceName: serviceName,
		logger:      slog.Default(), // используем стандартный логгер
	}
}

// GetLoggerFromCtx возвращает логгер из контекста
func GetLoggerFromCtx(ctx context.Context) Logger {
	if ctxLogger, ok := ctx.Value(LoggerKey).(Logger); ok {
		return ctxLogger
	}
	return New("default") // возвращаем логгер по умолчанию, если в контексте нет логгера
}