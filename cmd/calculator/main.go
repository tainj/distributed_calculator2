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

	example := calculator.NewExpression("2 + 3")
	fmt.Println(example.Convert())
	fmt.Println(example.Check())
	fmt.Println(example.Infix)
	fmt.Println(example.Postfix)
}