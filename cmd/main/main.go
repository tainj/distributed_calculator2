package main

import (
	"errors"
	"math"
)

var (
	ErrDivisionByZero = errors.New("You can't divide by zero.")
	ErrNonExistingOperation = errors.New("This operation does not exist or has not been implemented.")
)

type SimpleMathSolver interface {
	Calculate() (float64, error)
}

type SimpleExample struct {
	Num1 float64
	Num2 float64
	Sign string
}

func NewSimpleExample(num1, num2 float64, sign string) *SimpleExample {
	return &SimpleExample{Num1: num1, Num2: num2, Sign: sign}
}

func (s *SimpleExample) Calculate() (float64, error) {
	switch s.Sign {
	case "+":
		return s.Num1 + s.Num2, nil
	case "-":
		return s.Num1 - s.Num2, nil
	case "*":
		return s.Num1 * s.Num2, nil
	case "/":
		if s.Num2 == 0 {
			return 0, ErrDivisionByZero
		}
		return s.Num1 / s.Num2, nil
	case "^":
		return math.Pow(2, 3), nil
	default:
		return 0, ErrNonExistingOperation
	}
}

type Example struct {
	Example string
}

func (s *Example) Convert() () {

}

type 



func main() {
	s := "2 + 2"
}