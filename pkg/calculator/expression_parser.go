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

    // Операторы, которые правоассоциативны
    RightAssociative = map[string]bool{
        "^": true, // степень: 2^3^4 = 2^(3^4)
        // "~": true, // унарный минус — можно добавить, но он обрабатывается отдельно
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
            // обрабатываем операторы: +, -, *, /, ^, ~, (
            // выталкиваем из стека операторы с большим или равным приоритетом
            // НО: если оператор правоассоциативный (например, ^), то при равном приоритете — НЕ выталкиваем
            for !stack.IsEmptyStack() {
                top := stack.Peek()
                if top == "(" {
                    break
                }

                topPriority := OperatorPriority[top]

                // если приоритет верхнего больше текущего — выталкиваем
                if topPriority > value {
                    list = append(list, stack.Pop())
                } else if topPriority == value {
                    // если приоритет равен
                    // смотрим ассоциативность: если левоассоциативный — выталкиваем, если право — нет
                    if !RightAssociative[sign] {
                        list = append(list, stack.Pop())
                    } else {
                        break
                    }
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
            stack.Pop() // удаляем "("
        }
    }
    if len(number) > 0 {
        list = append(list, string(number)) // добавляем последние число, если имеется
    }
    for !stack.IsEmptyStack() {
        list = append(list, stack.Pop()) // выгружаем остаток стека
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
				var num1, num2 string
                var newExpr []string
                if sign == "~" {
                    // унарный минус: ~X → 0 - X
                    if index < 1 {
                        return nil, "" // ошибка: нет операнда
                    }
                    num1 = "0"
                    num2 = expression[index-1]
                    sign := "-" // всегда вычитание
                    result, variable := NewExample(num1, num2, sign)
                    results = append(results, result)
                    newExpr = replaceExpr(expression, index, variable)
                } else {
                    // Бинарный оператор: +, -, *, /, ^
                    if index < 2 {
                        return nil, "" // ошибка: мало операндов
                    }
                    num1 = expression[index-2]
                    num2 = expression[index-1]
                    result, variable := NewExample(num1, num2, sign)
                    results = append(results, result)
                    newExpr = replaceExpr(expression, index, variable)
                }

                expression = newExpr
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

