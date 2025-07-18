package calculator

import (
    "testing"
)

func TestConvert(t *testing.T) {
    tests := []struct {
        name     string // название теста
        input    string // входное выражение
        expected string // ожидаемый результат
        wantErr  bool   // ожидается ли ошибка
    }{
        {
            name:     "Simple addition",
            input:    "10 + 20",
            expected: "10 20 +",
            wantErr:  false,
        },
        {
            name:     "Simple subtraction",
            input:    "30 - 15",
            expected: "30 15 -",
            wantErr:  false,
        },
        {
            name:     "Simple multiplication",
            input:    "6 * 7",
            expected: "6 7 *",
            wantErr:  false,
        },
        {
            name:     "Simple division",
            input:    "100 / 20",
            expected: "100 20 /",
            wantErr:  false,
        },
        {
            name:     "Mixed operators without parentheses",
            input:    "10 + 20 * 30",
            expected: "10 20 30 * +",
            wantErr:  false,
        },
        {
            name:     "Mixed operators with parentheses",
            input:    "(10 + 20) * 30",
            expected: "10 20 + 30 *",
            wantErr:  false,
        },
        {
            name:     "Nested parentheses",
            input:    "(10 + (20 * 30)) / 5",
            expected: "10 20 30 * + 5 /",
            wantErr:  false,
        },
        {
            name:     "Complex expression",
            input:    "3 + 4 * 2 / (1 - 5)",
            expected: "3 4 2 * 1 5 - / +",
            wantErr:  false,
        },
        {
            name:     "Expression with multiple digits",
            input:    "100 + 25 * (3 - 4) + 5",
            expected: "100 25 3 4 - * + 5 +",
            wantErr:  false,
        },
        {
            name:     "Expression with same priority operators",
            input:    "10 - 5 + 3",
            expected: "10 5 - 3 +",
            wantErr:  false,
        },
        {
            name:     "Expression with decimal numbers",
            input:    "10.5 + 20.3",
            expected: "10.5 20.3 +",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            example := NewExample(tt.input)
            _, err := example.Convert()
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr && example.Postfix != tt.expected {
                t.Errorf("Convert() = %v, expected %v", example.Postfix, tt.expected)
            }
        })
    }
}