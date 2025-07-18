package calculator

import (
	"errors"
	"strings"
	"unicode"
)

var (
    ErrDivisionByZero       = errors.New("division by zero")
    ErrNonExistingOperation = errors.New("operation does not exist or not implemented")
    ErrCovertExample        = errors.New("line is not a mathematical expression or contains an error")
)

var (
	OperatorPriority = map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"^": 3,
		"~": 4,
		"(": 6,
	}
)
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
	return len(s.list) == 0
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

type Expression struct {
    Infix   string // инфиксное выражение
    Postfix string // постфиксное выражение
}

func NewExample(str string) *Expression {
	return &Expression{Infix: str}
}

func (s *Expression) Convert() (bool, error) {
	// инициализация списка, стека и списка для чисел
	list := make([]string, 0)
	stack := NewStack()
	example := strings.ReplaceAll(s.Infix, " ", "") // удаляем пробелы
	number := make([]rune, 0)
	for _, i := range example {
		sign := string(i)
		if unicode.IsDigit(i) { // проверяем является ли символ цифрой
			number = append(number, i) // добавляем в список для чисел
			continue
		} else if sign == "." {
			number = append(number, rune(sign[0]))
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
			stack.Push(sign) // добавляем текущий оператор в стек
		}
		if i == ')' { // извлекаем операторы из стека 
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
	s.Postfix = strings.Join(list, " ")
	return true, nil
}


