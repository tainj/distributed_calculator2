package main

import (
	"fmt"

	"github.com/tainj/distributed_calculator2/pkg/calculator"
	"github.com/tainj/distributed_calculator2/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	example := calculator.NewExpression("~2 + 3")
	example.Convert()
	fmt.Println(example.Postfix)
}