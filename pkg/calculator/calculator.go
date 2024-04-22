package calculator

import (
	"regexp"
	"strconv"
	"strings"
)

type Stack []string

// IsEmpty: check if stack is empty
func (st *Stack) IsEmpty() bool {
	return len(*st) == 0
}

// Push a new value onto the stack
func (st *Stack) Push(str string) {
	*st = append(*st, str) //Simply append the new value to the end of the stack
}

// Remove top element of stack. Return false if stack is empty.
func (st *Stack) Pop() bool {
	if st.IsEmpty() {
		return false
	} else {
		index := len(*st) - 1 // Get the index of top most element.
		*st = (*st)[:index]   // Remove it from the stack by slicing it off.
		return true
	}
}

// Return top element of stack. Return false if stack is empty.
func (st *Stack) Top() string {
	if st.IsEmpty() {
		return ""
	} else {
		index := len(*st) - 1   // Get the index of top most element.
		element := (*st)[index] // Index onto the slice and obtain the element.
		return element
	}
}

// Function to return precedence of operators
func prec(s string) int {
	if s == "^" {
		return 3
	} else if (s == "/") || (s == "*") {
		return 2
	} else if (s == "+") || (s == "-") {
		return 1
	} else {
		return -1
	}
}

func InfixToPostfix(infix string) string {
	var sta Stack
	var postfix string
	for _, char := range infix {
		opchar := string(char)
		// if scanned character is operand, add it to output string
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			postfix = postfix + opchar
		} else if char == '(' {
			sta.Push(opchar)
		} else if char == ')' {
			for sta.Top() != "(" {
				postfix = postfix + sta.Top()
				sta.Pop()
			}
			sta.Pop()
		} else {
			for !sta.IsEmpty() && prec(opchar) <= prec(sta.Top()) {
				postfix = postfix + sta.Top()
				sta.Pop()
			}
			sta.Push(opchar)
		}
	}
	// Pop all the remaining elements from the stack
	for !sta.IsEmpty() {
		postfix = postfix + sta.Top()
		sta.Pop()
	}
	return postfix
}

type Expression struct {
	Variables map[string]int
}

func ParseEquation(equation string) []string {
	re := regexp.MustCompile(`\d+|[()+*/-]`)
	elements := re.FindAllString(equation, -1)
	return elements
}

func Translate(expression string) (string, Expression, error) {
	ex := Expression{Variables: make(map[string]int)}
	alphabet := []string{"a", "b", "c", "d", "e", "f", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "x", "y", "z"}
	for _, i := range alphabet {
		ex.Variables[i] = 137  // нулевое значение
	}
	new_expression := []string{}
	for _, i := range ParseEquation(expression) {
		if !strings.Contains("()*+/-", i) {
			num, err := strconv.Atoi(i)
			if err != nil {
				return "", ex, err
			}
			alphabet2 := []string{}
			for val, _ := range ex.Variables {
				alphabet2 = append(alphabet2, val)
			}
			flag := ""
			for _, val := range alphabet {
				if ex.Variables[val] == 137 {
					ex.Variables[val] = num
					flag = val
					break
				}
			}
			new_expression = append(new_expression, flag)
		} else {
			new_expression = append(new_expression, i)
		}
	}
	return strings.Join(new_expression, ""), ex, nil
}

type Solution struct {
	Status int
	Result int
}

func Calc(ch chan Solution, sign string, a int, b int) {
	switch sign {
	case "+":
		ch <- Solution{Result: a + b, Status: 200}
	case "-":
		ch <- Solution{Result: a - b, Status: 200}
	case "*":
		ch <- Solution{Result: a * b, Status: 200}
	case "/":
		if b == 0 {
			ch <- Solution{Result: 0, Status: 456} // статус ошибки делении на 0
		} else {
			ch <- Solution{Result: a / b, Status: 200}
		}
	}
}

// func main() {
// 	expression := "2 - 2 * 6 + 7"
// 	s, ex, err := Translate(expression)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(s)
// 	fmt.Println(ex.variables["a"])
// 	infix := infixToPostfix(s)
// 	fmt.Println(infix)
// }
