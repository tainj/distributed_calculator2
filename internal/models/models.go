package models

import "github.com/tainj/distributed_calculator2/pkg/calculator"

type Example struct {
	Id             string
	Expression     string
	SimpleExamples []calculator.MathExample
	Response       string
}