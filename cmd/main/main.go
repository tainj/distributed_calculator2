package main

import (
	"fmt"
	"github.com/tainj/distributed_calculator2/pkg/calculator"
)

func main() {
	s := calculator.NewExample("100 + 25 * (3 - 4) + 5")
	fmt.Println(s.Check())
	s.Convert()
	fmt.Println(s.Calculate())
}