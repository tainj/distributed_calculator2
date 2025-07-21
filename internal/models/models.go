package models

import (
	"github.com/tainj/distributed_calculator2/pkg/calculator"
)

type Example struct {
	Id             string        `json:"id"`
	Expression     string        `json:"expression"`
	SimpleExamples []calculator.MathExample `json:"simple_examples"`
	Response       string        `json:"response"`
}