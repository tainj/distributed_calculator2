package middlewares

import (
    "context"
    "fmt"
    "net/http"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

// AuthMiddleware - HTTP middleware для аутентификации через JWT
func AuthMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Получаем заголовок Authorization
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing authorization header", http.StatusUnauthorized)
                return
            }

            // Проверяем формат Bearer
            if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
                http.Error(w, "invalid authorization format", http.StatusUnauthorized)
                return
            }

            // Извлекаем токен
            tokenStr := authHeader[7:]

            // Структура для claims (должна совпадать с gRPC)
            type jwtClaims struct {
                jwt.RegisteredClaims
                UserId string `json:"user_id"`
            }

            // Парсим и проверяем токен
            var claims jwtClaims
            token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
                // Проверяем метод подписи
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }
                // Секретный ключ (лучше брать из env)
                return []byte("твой_очень_секретный_ключ"), nil
            })

            // Обрабатываем ошибки
            if err != nil {
                http.Error(w, fmt.Sprintf("token parse error: %v", err), http.StatusUnauthorized)
                return
            }

            if !token.Valid {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }

            // Проверяем, что UserId - валидный UUID
            if _, err := uuid.Parse(claims.UserId); err != nil {
                http.Error(w, "invalid user ID format", http.StatusUnauthorized)
                return
            }

            // Создаем ключ для контекста (как в gRPC)
            type userIdCtxKey struct{}

            // Кладем userID в контекст
            ctx := context.WithValue(r.Context(), userIdCtxKey{}, claims.UserId)
            
            // Передаем запрос дальше с обновленным контекстом
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}