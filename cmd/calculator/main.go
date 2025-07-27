package main

import (
	"fmt"

	"github.com/tainj/distributed_calculator2/pkg/calculator"
)

func main() {
	// cfg, err := config.LoadConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(cfg)

	example := calculator.NewExpression("~(~5)")
	fmt.Println(example.Convert())
	fmt.Println(example.Infix)
	fmt.Println(example.Postfix)
}