package main

import (
	"github.com/tantoni228/distributed_calculator2/cmd/server"
)

func main() {
	server.CreateCalcServer(5002)
}