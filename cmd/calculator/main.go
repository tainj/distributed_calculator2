package main

import (
	"fmt"
	"github.com/tainj/distributed_calculator2/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.Grpc.RestPort)

}