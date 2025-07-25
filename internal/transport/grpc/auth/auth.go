package auth

import (
    "context"
    "fmt"
    "errors"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
)

// Ключ для контекста (должен быть экспортируемым или общим)
type UserIdCtxKey struct{}

// Claims для JWT
type JwtClaims struct {
    jwt.RegisteredClaims
    UserId string `json:"user_id"`
}

type SecurityHandler struct{}

func NewSecurityHandler() *SecurityHandler {
    return &SecurityHandler{}
}

func (s *SecurityHandler) AuthInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // Получаем metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }

    // Ищем заголовок authorization
    authHeaders, ok := md["authorization"]
    if !ok || len(authHeaders) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing authorization header")
    }

    authHeader := authHeaders[0]

    // Проверяем формат Bearer
    if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
        return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
    }
    tokenStr := authHeader[7:]

    // Парсим JWT токен
    var claims JwtClaims
    token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
        // Проверяем метод подписи
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        // Секретный ключ
        return []byte("твой_очень_секретный_ключ"), nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenUnverifiable) {
            return nil, status.Error(codes.Unauthenticated, "token verification failed")
        }
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }

    if !token.Valid {
        return nil, status.Error(codes.Unauthenticated, "token is invalid")
    }

    // Проверка UUID
    if _, err := uuid.Parse(claims.UserId); err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid user ID format")
    }

    // Кладем userID в контекст
    ctx = context.WithValue(ctx, UserIdCtxKey{}, claims.UserId)
    return handler(ctx, req)
}

// Получение userID из контекста
func UserIdFromCtx(ctx context.Context) string {
    if userId, ok := ctx.Value(UserIdCtxKey{}).(string); ok {
        return userId
    }
    return ""
}