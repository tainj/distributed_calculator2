package calculator

import (
	"reflect"
	"testing"
)

func TestNewExample_GeneratesVariable(t *testing.T) {
	example, variable := NewExample("2", "3", "+")
	if variable == "" {
		t.Error("NewExample() returned empty variable")
	}
	if example.Variable != variable {
		t.Error("Example.Variable != returned variable")
	}
	if example.Num1 != "2" || example.Num2 != "3" || example.Sign != "+" {
		t.Error("Example fields not set correctly")
	}
}

func TestNode_Calculate(t *testing.T) {
	tests := []struct {
		name     string
		num1     float64
		num2     float64
		sign     string
		expected float64
		wantErr  bool
	}{
		{"add", 2, 3, "+", 5, false},
		{"subtract", 5, 3, "-", 2, false},
		{"multiply", 4, 3, "*", 12, false},
		{"divide", 8, 4, "/", 2, false},
		{"divide by zero", 5, 0, "/", 0, true},
		{"invalid op", 2, 3, "%", 0, true},
		{"float", 2.5, 1.5, "+", 4.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Node{Num1: tt.num1, Num2: tt.num2, Sign: tt.sign}
			result, err := node.Calculate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("Calculate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewNode_GeneratesVariable(t *testing.T) {
	node:= NewNode(2, 3, "+")
	if node.Num1 != 2 || node.Num2 != 3 || node.Sign != "+" {
		t.Error("Node fields not set correctly")
	}
}

func TestReplaceExpr(t *testing.T) {
	tests := []struct {
		name     string
		expr     []string
		opIndex  int
		varName  string
		expected []string
	}{
		{
			name:     "basic replacement",
			expr:     []string{"2", "3", "+"},
			opIndex:  2,
			varName:  "tmp1",
			expected: []string{"tmp1"},
		},
		{
			name:     "middle of expression",
			expr:     []string{"tmp1", "4", "*"},
			opIndex:  2,
			varName:  "tmp2",
			expected: []string{"tmp2"},
		},
		{
			name:     "complex chain",
			expr:     []string{"2", "3", "+", "4", "*", "5", "+"},
			opIndex:  2,
			varName:  "tmp1",
			expected: []string{"tmp1", "4", "*", "5", "+"},
		},
		{
			name:     "out of bounds start",
			expr:     []string{"+", "3", "*"},
			opIndex:  0,
			varName:  "tmp",
			expected: []string{"tmp", "3", "*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceExpr(tt.expr, tt.opIndex, tt.varName)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("replaceExpr() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestExpression_Calculate(t *testing.T) {
	tests := []struct {
		name          string
		postfix       string
		expectedCount int
		finalIsVar    bool
	}{
		{
			name:          "Simple: 2 + 3",
			postfix:       "2 3 +",
			expectedCount: 1,
			finalIsVar:    true,
		},
		{
			name:          "Complex: 2 3 + 4 *",
			postfix:       "2 3 + 4 *",
			expectedCount: 2,
			finalIsVar:    true,
		},
		{
			name:          "Three operations",
			postfix:       "2 3 + 4 * 5 +",
			expectedCount: 3,
			finalIsVar:    true,
		},
		{
			name:          "Single number",
			postfix:       "42",
			expectedCount: 0,
			finalIsVar:    false, // final = "42"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &Expression{Postfix: tt.postfix}
			results, final := expr.Calculate()

			if len(results) != tt.expectedCount {
				t.Errorf("Calculate() returned %d tasks, expected %d", len(results), tt.expectedCount)
			}

			if tt.finalIsVar {
				// Проверим, что final совпадает с последней переменной
				if len(results) > 0 && final != results[len(results)-1].Variable {
					t.Errorf("Final variable = %s, expected last task var", final)
				}
			} else {
				// Должно быть число
				if final != tt.postfix {
					t.Errorf("Final = %s, expected %s", final, tt.postfix)
				}
			}
		})
	}
}

func TestExpression_Check(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid: 2+3", "2 + 3", true},
		{"valid: float", "2.5 * 4.1", true},
		{"valid: with parentheses", "(2 + 3) * 4", true},
		{"invalid: double op", "2 + + 3", false},
		{"invalid: unclosed paren", "(2 + 3", false},
		{"invalid: empty", "", false},
		{"valid: power", "2 ^ 3", true},
		{"invalid: trailing op", "2 +", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewExpression(tt.input)
			valid := expr.Check()
			if valid != tt.valid {
				t.Errorf("Check() = %v, want %v", valid, tt.valid)
			}
		})
	}
}

func TestExpression_Convert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple addition",
			input:    "10 + 20",
			expected: "10 20 +",
			wantErr:  false,
		},
		{
			name:     "Simple multiplication",
			input:    "6 * 7",
			expected: "6 7 *",
			wantErr:  false,
		},
		{
			name:     "With parentheses",
			input:    "(2 + 3) * 4",
			expected: "2 3 + 4 *",
			wantErr:  false,
		},
		{
			name:     "Nested parentheses",
			input:    "(10 + (20 * 30)) / 5",
			expected: "10 20 30 * + 5 /",
			wantErr:  false,
		},
		{
			name:     "Decimal numbers",
			input:    "2.5 + 3.7",
			expected: "2.5 3.7 +",
			wantErr:  false,
		},
		{
			name:     "Complex expression",
			input:    "3 + 4 * 2 / (1 - 5)",
			expected: "3 4 2 * 1 5 - / +",
			wantErr:  false,
		},
		{
			name:     "Right-to-left operator (^)",
			input:    "2 ^ 3 ^ 4",
			expected: "2 3 4 ^ ^",
			wantErr:  false,
		},
		{
			name:     "Invalid syntax",
			input:    "2 + + 3",
			expected: "",
			wantErr:  false, // govaluate может пропустить, но Check() поймает
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Single number",
			input:    "42",
			expected: "42",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := NewExpression(tt.input)
			_, err := expr.Convert()

			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && expr.Postfix != tt.expected {
				t.Errorf("Convert() postfix = %q, expected %q", expr.Postfix, tt.expected)
			}
		})
	}
}

func TestStack(t *testing.T) {
	tests := []struct {
		name     string
		ops      []struct{ op string; val string }
		expected []string
	}{
		{
			name: "push and pop",
			ops: []struct{ op string; val string }{
				{"push", "a"},
				{"push", "b"},
				{"pop", "b"},
				{"pop", "a"},
			},
			expected: nil,
		},
		{
			name: "peek",
			ops: []struct{ op string; val string }{
				{"push", "x"},
				{"push", "y"},
				{"peek", "y"},
				{"pop", "y"},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewStack()
			for _, op := range tt.ops {
				switch op.op {
				case "push":
					stack.Push(op.val)
				case "pop":
					got := stack.Pop()
					if got != op.val {
						t.Errorf("Pop() = %v, want %v", got, op.val)
					}
				case "peek":
					got := stack.Peek()
					if got != op.val {
						t.Errorf("Peek() = %v, want %v", got, op.val)
					}
				}
			}
		})
	}
}