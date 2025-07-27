package auth

import "golang.org/x/crypto/bcrypt"


// HashPassword хэширует пароль
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}


// CheckPassword сравнивает открытый пароль с хэшированным
func CheckPassword(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}