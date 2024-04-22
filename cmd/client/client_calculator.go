package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	// "github.com/tantoni228/distributed_calculator2/cmd/server"
	"github.com/tantoni228/distributed_calculator2/cmd/server"
	"github.com/tantoni228/distributed_calculator2/database/connection"
	"github.com/tantoni228/distributed_calculator2/pkg/calculator"
	pb "github.com/tantoni228/distributed_calculator2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

type Data struct {
    Error string
	Id  int
	Solution []string
}

type Examples struct {
	Variables map[string]int
}

var wg sync.WaitGroup

var map_status_servers = make(map[int]int)


func handlerCalc(w http.ResponseWriter, r *http.Request) {
	data := Data{
	 Error: "Нет ошибок",
	}
	fmt.Println(r.Method)
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("front/index.html")
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}
	
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, data)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		token := r.FormValue("token")
		example := r.FormValue("example")
        const hmacSampleSecret = "yandex_luceum"
        tokenFromString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
            }
    
            return []byte(hmacSampleSecret), nil
        })
		log.Printf("%v\n", token)
    
        if err != nil {
            data = Data{
                Error: fmt.Sprintf("%v", err),
               }
        }
    
        if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
            fmt.Println("user name: ", claims["name"])
			num, err := connection.AddExample(example, fmt.Sprintf("%v", claims["name"]))
			if err != nil {
				data = Data{
					Error: fmt.Sprintf("%v", err),
				   }
			} else {
				data = Data{
					Error: "Нет ошибок",
					Id: int(num),
				}
			}
        } else {
            data = Data{
                Error: fmt.Sprintf("%v", err),
               }
        }

		
        tmpl, err := template.ParseFiles("front/index.html")
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}
	
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, data)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetRequest(port string, a int, b int, sign string, letter string) (int, error) {
	host := "localhost"

	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()
	/// ..будет продолжение

	grpcClient := pb.NewDistributedCalculatorClient(conn)

	example, err := grpcClient.Calculation(context.TODO(), &pb.SimpleExpression{
		Sign:   sign,
		A:      int32(a),
		B:      int32(b),
		Letter: letter,
	})

	if err != nil {
		log.Println("failed invoking Calculation: ", err)
		return 0, err
	}

	if example.Status != 200 {
		switch int(example.Status) {
		case 456:
			log.Println("Division by zero is not possible")
			return 0, fmt.Errorf("Division by zero is not possible")
		}
	}

	fmt.Println("Example: ", example.Result)
	return int(example.Result), nil
}

func Calculator2(infix string, ex calculator.Expression, port string) (string, calculator.Expression, []string, error) {
	infix_rune := []rune(infix)
	s := "" // сверху описано
	g := "" // для хранения переменой замены двух переменых и знака операции
	data := []string{}
	for val, sign := range infix_rune {
		if val + 1 > (len(infix_rune) - 2) {
			continue
		}
		alphabet := strings.Join([]string{"a", "b", "c", "d", "e", "f", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "x", "y", "z"}, "")
		// fmt.Println(string(sign), string(infix_rune[val+1]), string(infix_rune[val+2]))
		if strings.Contains(alphabet, string(sign)) && strings.Contains(alphabet, string(infix_rune[val+1])) && strings.Contains("+-/*", string(infix_rune[val+2])) {
			a := string(sign)                                     // переменая
			b := string(infix_rune[val+1])                        // переменая
			signCalc := string(infix_rune[val+2])                 // знак
			s = string(a) + string(b) + string(signCalc)          // соединяю чтобы заменить на переменую
			fmt.Println(s)
			if ex.Variables[a] != 137 && ex.Variables[b] != 137 { // ищем пустую переменую
				for i, val := range ex.Variables {
					if val == 137 {
						g = i
						fmt.Println(i)
						break
					}
				}
			}
			fmt.Println(ex.Variables[a], ex.Variables[b], signCalc)
			num, err := GetRequest(port, ex.Variables[a], ex.Variables[b], signCalc, g)
			if err != nil {
				return infix, ex, data, err
			}
			data = append(data, fmt.Sprintf("%v %v %v = %v\n", ex.Variables[a], signCalc, ex.Variables[b], num))
			ex.Variables[g] = num
			infix = strings.Replace(infix, s, g, -1)
			fmt.Println(infix)
			for letter, value := range ex.Variables {
				if value != 137 {
					fmt.Printf("%s: %d\n", letter, value)
				}
			}
		}
	}
	return infix, ex, data, nil
}

func Calculator() {
	example, err := connection.GetExampleByInfixLength()
	if err != nil {
		if err == sql.ErrNoRows {
			return 
		}
		// ошибка
		return 
	}

	// fmt.Println(example.Variables)
	var variables map[string]int
	var hist connection.Hist
	if err := json.Unmarshal([]byte(example.Variables), &variables); err != nil {
		fmt.Println("Что то не так")
		return 
	}

	if err := json.Unmarshal([]byte(example.History), &hist); err != nil {
		fmt.Println("Что то не так")
		return 
	}

	fmt.Println(hist)


	// fmt.Println(variables)
	ex := calculator.Expression{Variables: variables}
	fmt.Println(ex)
	infix := example.Infix


	if len(infix) == 1 {
		for letter, value := range ex.Variables {
			if value != 137 {
				fmt.Printf("%s: %d\n", letter, value)
			}
		}
		// fmt.Println(infix)
		// fmt.Println("Solution: ", ex.Variables[infix])
	} else {
		log.Printf("start calculatin %v\n", example.Id)
		fmt.Printf("infix: %v ex: %v\n", infix, ex)
		infix, ex, data, err := Calculator2(infix, ex, "5002")
		if err != nil {
			fmt.Println(err)
		}

		for _, i := range data {
			hist.History = append(hist.History, i)
		}

		fmt.Println(hist)

		err = connection.UpdateExample(int(example.Id), infix, ex, hist)
		if err != nil {
			panic(err)
		}
		log.Println("sucsefully change")
	}

}

func handleGetSolution(w http.ResponseWriter, r *http.Request) {
		data := Data{
	 Error: "Нет ошибок",
	}
	fmt.Println(r.Method)
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("front/solution_form.html")
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}
	
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, data)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		token := r.FormValue("token")
		id := r.FormValue("id")
        const hmacSampleSecret = "yandex_luceum"
        tokenFromString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                log.Printf(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
            }
            return []byte(hmacSampleSecret), nil
        })
		log.Printf("%v\n", token)
    
        if err != nil {
            data = Data{
                Error: fmt.Sprintf("%v", err),
               }
        }

		num, err := strconv.Atoi(id)
		if err != nil {
			data = Data{
                Error: fmt.Sprintf("%v", err),
               }
		} else {
			if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
				fmt.Println("user name: ", claims["name"])
				ok, err := connection.GetUserID(fmt.Sprintf("%v", claims["name"]), num)
				if err != nil {
					data = Data{
						Error: fmt.Sprintf("%v", err),
					   }
				}
				if !ok && err == nil {
					data = Data{
						Error: fmt.Sprint("Вы не имеете доступ к данному примеру."),
					   }
				}
				hist, err := connection.GetHistoryAsJSON(num)
				if err != nil {
					data = Data{
						Error: fmt.Sprintf("%v", err),
					   }
				} else {
					data = Data{
						Error: fmt.Sprint("Нет ошибок"),
						Solution: hist.History,
					}
				}
			} else {
				data = Data{
					Error: fmt.Sprintf("%v", err),
				   }
			}
		}
        

		
        tmpl, err := template.ParseFiles("front/solution_form.html")
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}
	
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, data)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}



func main() {

	wg.Add(3)
	go func() {
		go server.CreateCalcServer(5002)
	    map_status_servers[5002] = 200
	    wg.Done()
	}()
	go func(){
		defer wg.Done()
		log.Println("start change")
		for {
			time.Sleep(time.Second) 
			Calculator()
		}
		
	}()
	go func() {
		defer wg.Done()
		http.HandleFunc("/calculate", handlerCalc)
		http.HandleFunc("/get_solution", handleGetSolution)
        http.ListenAndServe(":8080", nil)
	}()
	wg.Wait()
}