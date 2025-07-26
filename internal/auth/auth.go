// internal/auth/auth_service.go
package auth

import (
    "errors"
    "time"
    
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
    
    "github.com/tainj/distributed_calculator2/internal/models"
)

type AuthService struct {
    userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
    return &AuthService{
        userRepo: userRepo,
    }
}

// HashPassword хеширует пароль
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// CheckPassword проверяет пароль
func CheckPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(username, email, password string) (*models.User, error) {
    // Проверяем, существует ли пользователь
    if _, err := s.userRepo.GetByUsername(username); err == nil {
        return nil, errors.New("username already exists")
    }
    
    if _, err := s.userRepo.GetByEmail(email); err == nil {
        return nil, errors.New("email already exists")
    }
    
    // Хешируем пароль
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return nil, err
    }
    
    // Создаем пользователя
    user := &models.User{
        ID:           uuid.New().String(),
        Username:     username,
        Email:        email,
        PasswordHash: hashedPassword,
        Role:         models.UserRole,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    // Сохраняем в БД
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

// Login авторизует пользователя
func (s *AuthService) Login(username, password string) (*models.User, error) {
    // Получаем пользователя
    user, err := s.userRepo.GetByUsername(username)
    if err != nil {
        // Пробуем найти по email
        user, err = s.userRepo.GetByEmail(username)
        if err != nil {
            return nil, errors.New("invalid credentials")
        }
    }
    
    // Проверяем пароль
    if !CheckPassword(user.PasswordHash, password) {
        return nil, errors.New("invalid credentials")
    }
    
    return user, nil
}