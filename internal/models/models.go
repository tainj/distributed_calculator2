package models

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