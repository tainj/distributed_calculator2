package models

import "time"

// для бизнес логики
type Example struct {
	Id             string        `json:"id"`
	Expression     string        `json:"expression"`
	SimpleExamples []Task        `json:"simple_examples"`
	Response       string        `json:"response"`
}

type Task struct {
    Num1     string `json:"num1"`
    Num2     string `json:"num2"`
    Sign     string `json:"sign"`
    Variable string `json:"variable"`
}

// для бд

type UserExample struct {
    ID         string     `json:"id" gorm:"primaryKey"`
    Expression string     `json:"expression"`
    Response   string     `json:"response"`   // variable финального результата
    Calculated bool       `json:"calculated"`
    Result     *float64   `json:"result,omitempty"` // может быть nil
    UserID     string     `json:"user_id"`    // для авторизации
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
}

type Step struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    ExampleID string    `json:"example_id"`
    Value1    float64   `json:"value1"`
    Value2    float64   `json:"value2"`
    Result    float64   `json:"result"`
    Sign      string    `json:"sign"`
    Variable  string    `json:"variable"`
    Order     int       `json:"order"`
    CreatedAt time.Time `json:"created_at"`
}