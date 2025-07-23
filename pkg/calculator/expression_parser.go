package calculator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/tainj/distributed_calculator2/internal/models"
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

func NewExample(num1, num2, sign string) (models.Task, string) {
	variable := uuid.New().String()  // генерируем имя переменной, куда будет сохранен результат
	return models.Task{Num1: num1, Num2: num2, Sign: sign, Variable: variable}, variable
}

// реализация стека и его методов
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

func NewExpression(str string) *Expression {
	return &Expression{Infix: str}
}

func (s *Expression) Check() bool {  // проверяется выражение на то, является ли оно корректным
	_, err := govaluate.NewEvaluableExpression(s.Infix)
	return err == nil
}

func (s *Expression) Convert() (bool, error) {
	if !s.Check() {
		return false, fmt.Errorf("line is not a mathematical expression or contains an error")
	}
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

func (s *Expression) Calculate() ([]models.Task, string) {
	results := make([]models.Task, 0)
	expression := strings.Split(s.Postfix, " ")  // формируем список из чисел и операторов
	for len(expression) != 1 {
		for index, sign := range expression {
			if _, isOperator := OperatorPriority[sign]; isOperator {
				example := expression[index - 2:index + 1]  // получаем два числа и оператор
				num1 := example[0]
				num2 := example[1]
				sign := example[2]
				result, variable := NewExample(num1, num2, sign)  // формируем пример
				results = append(results, result)
				expression = replaceExpr(expression, index, variable)
				break
			}
		}
	}
	return results, expression[0]
}

func replaceExpr(expr []string, opIndex int, varName string) []string {
    // определяем границы для вырезания
    start := opIndex - 2
    if start < 0 {
        start = 0 // чтоб не уйти в минус
    }

    end := opIndex + 1
    if end > len(expr) {
        end = len(expr) // чтоб не выйти за пределы
    }

    // собираем новый слайс
    newExpr := make([]string, 0, len(expr)-2)
    newExpr = append(newExpr, expr[:start]...) // всё до оператора
    newExpr = append(newExpr, varName)         // наша переменная
    newExpr = append(newExpr, expr[end:]...)   // всё после оператора

    return newExpr
}

