package main

import (
	"fmt"
	"os"

	"github.com/tainj/distributed_calculator2/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(cfg.Postgres.Host)
}