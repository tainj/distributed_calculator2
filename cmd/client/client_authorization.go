package main

import (
	"fmt"
	"log"
	"text/template"
	"context"
	"net/http"
	"os"
	"github.com/tantoni228/distributed_calculator2/cmd/server"
	pb "github.com/tantoni228/distributed_calculator2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

type Data struct {
	Error string
	Token string
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := Data{
	 Error: "Нет ошибок",
	}
	fmt.Println(r.Method)
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("front/register_form.html")
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
		login := r.FormValue("login")
		password := r.FormValue("password")
			  // сервер для авторизации будет на порту 5000
		host := "localhost"
		port := "5000"

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

		grpcClient := pb.NewAuthenticationServerClient(conn)

		auth, err := grpcClient.Authorization(context.TODO(), &pb.Form{
			Login: login,
			Password:  password,
		})

		if err != nil {
			log.Println("failed invoking Authorization: ", err)
		}
		switch (auth.Status) {
		case 200:
			data = Data{
				Error: "Нет ошибок",
			}
		case 261:
			data = Data{
				Error: "Пользователь с таким логином уже есть.",
			}
		default:
			data = Data{
				Error: "Возникли ошибки при авторизации.",
			}
		}
		tmpl, err := template.ParseFiles("front/register_form.html")
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

func handler2(w http.ResponseWriter, r *http.Request) {
	data := Data{
	 Error: "Нет ошибок",
	}
	fmt.Println(r.Method)
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("front/login_form.html")
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
		login := r.FormValue("login")
		password := r.FormValue("password")
			  // сервер для авторизации будет на порту 5000
		host := "localhost"
		port := "5000"

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

		grpcClient := pb.NewAuthenticationServerClient(conn)

		auth, err := grpcClient.Authentication(context.TODO(), &pb.Form{
			Login: login,
			Password:  password,
		})

		if err != nil {
			log.Println("failed invoking Authorization: ", err)
		}
		switch (auth.Status) {
		case 200:
			data = Data{
				Error: "Нет ошибок",
				Token: auth.Token,
			}
		case 261:
			data = Data{
				Error: "Пользователь с таким логином не найден в базе данных.",
				Token: auth.Token,
			}
		case 371:
			data = Data{
				Error: "Неверный пароль.",
				Token: auth.Token,
			}
		default:
			data = Data{
				Error: "Возникли ошибки.",
				Token: auth.Token,
			}
		}
		tmpl, err := template.ParseFiles("front/login_form.html")
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
	go server.CreateAuthServer(5000)

	http.HandleFunc("/register", handler)
	http.HandleFunc("/login", handler2)
    http.ListenAndServe(":5001", nil)
}