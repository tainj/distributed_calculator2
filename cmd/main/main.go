package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode"
)

var (
	ErrDivisionByZero       = errors.New("You can't divide by zero.")
	ErrNonExistingOperation = errors.New("This operation does not exist or has not been implemented.")
	ErrCovertExample        = errors.New("This line is not a mathematical expression or contains an error.")
)

var (
	OperatorPriority = map[string]int{
		"(": 6,
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"^": 3,
		"~": 4,
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
	list []string
}

func NewStack() *Stack {
	return &Stack{list: make([]string, 0)}
}

func (s *Stack) Push(item string) {
	s.list = append(s.list, item)
}

func (s *Stack) IsEmptyStack() bool {
	if len(s.list) == 0 {
		return true
	}
	return false
}

func (s *Stack) Pop() string {
	index := len(s.list) - 1
	result := s.list[index]
	s.list = s.list[:index]
	return result
}

func (s *Stack) Peek() string {
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
	// инициализация списка, стека и списка для чисел
	list := make([]string, 0)
	stack := NewStack()
	example := strings.ReplaceAll(s.InfixExpr, " ", "") // удаляем пробелы
	number := make([]rune, 0)
	for _, i := range example {
		sign := string(i)
		if unicode.IsDigit(i) { // проверяем является ли символ цифрой
			number = append(number, i) // добавляем в список для чисел
			continue
		} else {
			if len(number) != 0 {
				list = append(list, string(number)) // если это не цифра, то добавляем всю строку в список
				number = make([]rune, 0)
			}
		}
		if value, isOperator := OperatorPriority[sign]; isOperator {
			for !stack.IsEmptyStack() {
				top := stack.Peek()
				if top != "(" && OperatorPriority[top] >= value {
					list = append(list, stack.Pop()) // извлекаем оператор из стека
				} else {
					break
				}
			}
			stack.Push(sign) // Добавляем текущий оператор в стек
		}
		if i == ')' {
			for stack.Peek() != "(" {
				list = append(list, stack.Pop())
			}
			stack.Pop()
		}
	}
	if len(number) > 0 {
		list = append(list, string(number)) // добавляем последние число, если имеется
	}
	for !stack.IsEmptyStack() {
		list = append(list, stack.Pop())
	}
	return strings.Join(list, " "), nil
}

func main() {
	// s := NewExample("3 + 4 * 2 / (1 - 5)")
	// fmt.Println(s.Convert())
	// s2 := NewExample("10 + 25 * (3 - 4)")
	// fmt.Println(s2.Convert())
	s3 := NewExample("100 + 25 * (3 - 4) + 5")
	fmt.Println(s3.Convert())
}
