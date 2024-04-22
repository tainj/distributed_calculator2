package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"github.com/tantoni228/distributed_calculator2/cmd/server"
	"github.com/tantoni228/distributed_calculator2/pkg/calculator"
	pb "github.com/tantoni228/distributed_calculator2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

var wg sync.WaitGroup

var map_status_servers = make(map[int]int)

// 200 - живой; 467 - работает; 331 - умер


func main() {
	wg.Add(1)
	// go server.CreateCalcServer(5001)
	// map_status_servers[5002] = 200
	// wg.Done()
	go func() {
		go server.CreateCalcServer(5002)
	    map_status_servers[5002] = 200
	    wg.Done()
	}()
	wg.Wait()
	expression := "2 + 2"

	expression, ex, err := calculator.Translate(expression)
	if err != nil {
		fmt.Println(err)
		return
	}

	infix := calculator.InfixToPostfix(expression)

	fmt.Println(infix)
	// for letter, value := range ex.Variables {
	// 	fmt.Printf("%s: %d\n", letter, value)
	// }
	// fmt.Println(infix)

	for {
		if len(infix) == 1 {
			for letter, value := range ex.Variables {
				if value != 137 {
					fmt.Printf("%s: %d\n", letter, value)
				}
			}
			fmt.Println(infix)
			fmt.Println("Solution: ", ex.Variables[infix])
			break
		} else {
			infix, ex, err = Calculator2(infix, ex, "5002")
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
	

}
