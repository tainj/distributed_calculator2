package models

import "time"

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

// для бизнес логики

type Example struct {
	Id             string `json:"id"`
	Expression     string `json:"expression"`
	SimpleExamples []Task `json:"simple_examples"`
	Response       string `json:"response"`
}

type Task struct {
	Num1      string `json:"num1"`
	Num2      string `json:"num2"`
	Sign      string `json:"sign"`
	Variable  string `json:"variable"`
	ExampleID string `json:"example_id"`
	Index     int    `json:"index"`
	IsFinal   bool   `json:"is_final"`
}

// для бд

type UserExample struct {
	ID         string   `json:"id" db:"id"`
	Expression string   `json:"expression" db:"expression"`
	Response   string   `json:"response" db:"response"`
	Calculated bool     `json:"calculated" db:"calculated"`
	Result     *float64 `json:"result,omitempty" db:"result"`
	Error      *string  `json:"error,omitempty" db:"error"`
	UserID     string   `json:"user_id" db:"user_id"`
}

type Step struct {
	ID        string  `json:"id" db:"id"`
	ExampleID string  `json:"example_id" db:"example_id"`
	Value1    float64 `json:"value1" db:"value1"`
	Value2    float64 `json:"value2" db:"value2"`
	Result    float64 `json:"result" db:"result"`
	Sign      string  `json:"sign" db:"sign"`
	Variable  string  `json:"variable" db:"variable"`
	Order     int     `json:"order" db:"order"`
}

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	Role         Role      `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}