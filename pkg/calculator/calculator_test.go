package calculator

import (
    "testing"
)

func TestConvert(t *testing.T) {
    tests := []struct {
        name     string // название теста
        input    string // входное выражение
        expected string // ожидаемый результат
    }{
        {
            name:     "Simple addition",
            input:    "10 + 20",
            expected: "10 20 +",
        },
        {
            name:     "Simple subtraction",
            input:    "30 - 15",
            expected: "30 15 -",
        },
        {
            name:     "Simple multiplication",
            input:    "6 * 7",
            expected: "6 7 *",
        },
        {
            name:     "Simple division",
            input:    "100 / 20",
            expected: "100 20 /",
        },
        {
            name:     "Mixed operators without parentheses",
            input:    "10 + 20 * 30",
            expected: "10 20 30 * +",
        },
        {
            name:     "Mixed operators with parentheses",
            input:    "(10 + 20) * 30",
            expected: "10 20 + 30 *",
        },
        {
            name:     "Nested parentheses",
            input:    "(10 + (20 * 30)) / 5",
            expected: "10 20 30 * + 5 /",
        },
        {
            name:     "Complex expression",
            input:    "3 + 4 * 2 / (1 - 5)",
            expected: "3 4 2 * 1 5 - / +",
        },
        {
			name:     "Expression with multiple digits",
			input:    "100 + 25 * (3 - 4) + 5",
			expected: "100 25 3 4 - * + 5 +",
		},
        {
            name:     "Expression with same priority operators",
            input:    "10 - 5 + 3",
            expected: "10 5 - 3 +",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            example := NewExample(tt.input)
            result, err := example.Convert()
            if err != nil {
                t.Errorf("Convert() returned an error: %v", err)
            }
            if result != tt.expected {
                t.Errorf("Convert() = %v, expected %v", result, tt.expected)
            }
        })
    }
}