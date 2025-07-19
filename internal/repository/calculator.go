package repository

import "log"

type CalculatorRepository interface {
	SaveResult(expression string, result float64) error
	GetResult(expression string) (float64, error)
}

type calculatorRepository struct {
}

func NewCalculatorRepository() CalculatorRepository {
	return &calculatorRepository{}
}

func (r *calculatorRepository) SaveResult(expression string, result float64) error {
	// Заглушка: просто логируем
	log.Printf("Saved result: %s = %v", expression, result)
	return nil
}

func (r *calculatorRepository) GetResult(expression string) (float64, error) {
	// Заглушка: возвращаем любое значение
	return 42, nil
}