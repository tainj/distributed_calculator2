package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode"
)

var (
	ErrDivisionByZero = errors.New("You can't divide by zero.")
	ErrNonExistingOperation = errors.New("This operation does not exist or has not been implemented.")
	ErrCovertExample = errors.New("This line is not a mathematical expression or contains an error.")
)

var (
	OperatorPriority = map[rune]int{
		'(': 0,
        '+': 1,
        '-': 1,
        '*': 2,
        '/': 2,
        '^': 3,
        '~': 4,
	}
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

type Stack struct {
	list []rune
}

func NewStack() *Stack {
	return &Stack{list: make([]rune, 0)}
}

func (s *Stack) Push(item rune) {
	s.list = append(s.list, item)
}

func (s *Stack) IsEmptyStack() bool {
	if len(s.list) == 0 {
		return true
	}
	return false
}

func (s *Stack) Pop() rune {
	index := len(s.list) - 1
	result := s.list[index]
	s.list = s.list[:index]
	return result
}

func (s *Stack) Peek() rune {
	index := len(s.list) - 1
	return s.list[index]
}

type Example struct {
	InfixExpr string
	// PostfixExp string
}

func NewExample(str string) *Example {
	return &Example{InfixExpr: str}
}


func (s *Example) Convert() (string, error) {
	list := make([]rune, 0)
	stack := NewStack()
	example := strings.ReplaceAll(s.InfixExpr, " ", "")
	for _, i := range example {
		if unicode.IsDigit(i) {
			list = append(list, i)
		}
		for key, value := range OperatorPriority {
			if key == i {
				if stack.IsEmptyStack() {
					stack.Push(i)
					continue
				}
				if OperatorPriority[i] > OperatorPriority[stack.Peek()] {
					stack.Push(i)
					continue
				}
				if OperatorPriority[i] == value {
					list = append(list, stack.Pop())
					stack.Push(i)
				}
			}
		}
		if i == ')' {
			for stack.Peek() != '(' {
				list = append(list, stack.Pop())
			}
		}
	}
	for !stack.IsEmptyStack() {
		list = append(list, stack.Pop())
	}
	fmt.Println(list)
	return string(list), nil
}


func main() {
	s := NewExample("3 + 4 * 2 / (1 - 5)")
	fmt.Println(s.Convert())
}