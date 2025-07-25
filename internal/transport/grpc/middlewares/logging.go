package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tainj/distributed_calculator2/internal/transport/grpc/auth"
	"github.com/tainj/distributed_calculator2/pkg/logger"
)

// LoggerProvider добавляет логгер в контекст запроса
func LoggerProvider(serviceName string) Middleware {
    l := logger.New(serviceName)
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), logger.LoggerKey, l)
            r = r.WithContext(ctx)
            next.ServeHTTP(w, r)
        })
    }
}

// loggingResponseWriter — для отслеживания статус-кода
type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
    return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}

// Logging middleware — логирует завершённые запросы
func Logging() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            l := logger.GetLoggerFromCtx(r.Context())
            lrw := newLoggingResponseWriter(w)
            start := time.Now()

            defer func() {
                duration := time.Since(start).Milliseconds()
                userId := auth.UserIdFromCtx(r.Context()) // Получаем userId из контекста

                // Формируем атрибуты для лога
                attrs := []any{
                    "method", r.Method,
                    "uri", r.RequestURI,
                    "status_code", lrw.statusCode,
                    "elapsed_ms", duration,
                }
                
                // Добавляем userId только если он есть
                if userId != "" {
                    attrs = append(attrs, "user_id", userId)
                }

                l.Info(r.Context(),
                    fmt.Sprintf("%s request to %s completed", r.Method, r.RequestURI),
                    attrs...,
                )
            }()

            next.ServeHTTP(lrw, r)
        })
    }
}