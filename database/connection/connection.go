package connection

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tantoni228/distributed_calculator2/pkg/calculator"
)

func createTables(ctx context.Context, db *sql.DB) error {
	queries := []string{
        `
        CREATE TABLE IF NOT EXISTS users(
            id INTEGER PRIMARY KEY AUTOINCREMENT, 
            login TEXT UNIQUE,
            password TEXT
        );`,
        `
        CREATE TABLE IF NOT EXISTS example(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            infix TEXT,
            variables JSON,
            id_user INTEGER,
            history JSON,
			error TEXT,
            FOREIGN KEY (id_user) REFERENCES users(id)
        );`,
    }

    for _, query := range queries {
        if _, err := db.ExecContext(ctx, query); err != nil {
            return err
        }
    }

    return nil
}

type Hist struct {
    History []string
}

func InsertUser(login, password string) error {
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "database/calculator.db")
	if err != nil {
		log.Println("open db error")
		panic(err)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		log.Println("connect error")
		panic(err)
	}

	if err = createTables(ctx, db); err != nil {
		log.Println("create table error")
		panic(err)
	}

	const insertUserQuery = `
	  INSERT INTO users (login, password) VALUES(?, ?);
	`
	
	_, err = db.ExecContext(ctx, insertUserQuery, login, password)
	if err != nil {
		log.Printf("user entry error")
		return err
	}
	return nil
}

func GenerateTokenUser(login, password string) (string, error){
	var need_password string
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "database/calculator.db")
	if err != nil {
		log.Println("open db error")
		panic(err)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		log.Println("connect error")
		panic(err)
	}

	if err = createTables(ctx, db); err != nil {
		log.Println("create table error")
		panic(err)
	}
	err = db.QueryRow("SELECT login, password FROM users WHERE login=?", login).Scan(&login, &need_password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("the user with this username was not found")
			return "", err
		}
		log.Fatal(err)
	}
	if password != need_password {
		log.Println("invalid password")
			return "", fmt.Errorf("invalid password")
	}

	const hmacSampleSecret = "yandex_luceum"
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": login,
		"nbf":  now.Unix(),
		"exp":  now.Add(15 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		panic(err)
	}

	return tokenString, nil
}

func getIdUser(ctx context.Context, db *sql.DB, login string) (int, error){
	var id int
	query := `
		SELECT id
		FROM users
		WHERE login = ?;`
	err := db.QueryRow(query, login).Scan(&id)
	if err != nil {
		log.Fatal(err)
		if err != nil {
			return 0, err
		}
	}
	return id, err
}

func AddExample(expression string, user string) (int64, error) {

    // подключение к бд
    ctx := context.TODO()
    db, err := sql.Open("sqlite3", "database/calculator.db")
    if err != nil {
        log.Println("open db error")
        panic(err)
    }
    defer db.Close()

    err = db.PingContext(ctx)
    if err != nil {
        log.Println("connect error")
        panic(err)
    }

    if err = createTables(ctx, db); err != nil {
        log.Println("create table error")
        panic(err)
    }

    // переводим выражение
    expression, ex, err := calculator.Translate(expression)
    if err != nil {
        fmt.Println(err)
        return 0, err
    }

    expression = calculator.InfixToPostfix(expression)

    jsonBytes, err := json.Marshal(ex.Variables)
    if err != nil {
        return 0, err
    }

    history := map[string]interface{}{"History": []string{"Solution\n"}}
    jsonBytes2, err := json.Marshal(history)
    if err != nil {
        return 0, err
    }

    id, err := getIdUser(ctx, db, user)
    if err != nil {
        return 0, err
    }

    const insertUserQuery = `
        INSERT INTO example (variables, id_user, history, infix) VALUES(?, ?, ?, ?) RETURNING id;
    `
    fmt.Println(string(jsonBytes), string(jsonBytes2))
    var lastId int64
    row := db.QueryRowContext(ctx, insertUserQuery, string(jsonBytes), id, jsonBytes2, expression)
    err = row.Scan(&lastId)
    if err != nil {
        log.Printf("user entry error")
        return 0, err
    }
    return lastId, nil
}

type Examples struct {
	Variables map[string]int
	Expression string
}

type Example struct {
	Id          int64           `json:"id"`
	Infix       string          `json:"infix"`
	Variables   string          `json:"variables"`
	IDUser      int             `json:"id_user"`
	History     string          `json:"history"`
}


func GetExampleByInfixLength() (Example, error) {
	ctx := context.TODO()
    db, err := sql.Open("sqlite3", "database/calculator.db")
    if err != nil {
        log.Println("open db error")
        panic(err)
    }
    defer db.Close()

    err = db.PingContext(ctx)
    if err != nil {
        log.Println("connect error")
        panic(err)
    }

    if err = createTables(ctx, db); err != nil {
        log.Println("create table error")
        panic(err)
    }
    // Выбираем одну запись из таблицы.
    row := db.QueryRowContext(ctx, `
        SELECT id, infix, variables, id_user, history
        FROM example
        WHERE LENGTH(infix) > 1;`)

    // Создаём объект для хранения данных.\


    var example Example

    // Сканируем данные из строки в объект.
    if err := row.Scan(&example.Id, &example.Infix, &example.Variables, &example.IDUser, &example.History); err != nil {
        fmt.Println(err)
        return Example{}, err
    }
    // fmt.Println(err)
    fmt.Println(example)
    // Возвращаем объект.
    return example, nil
}

func UpdateExample(id int, infix string, ex calculator.Expression, hist Hist) error {
	ctx := context.TODO()
    db, err := sql.Open("sqlite3", "database/calculator.db")
    if err != nil {
        log.Println("open db error")
        panic(err)
    }
    defer db.Close()

    err = db.PingContext(ctx)
    if err != nil {
        log.Println("connect error")
        panic(err)
    }

    if err = createTables(ctx, db); err != nil {
        log.Println("create table error")
        panic(err)
    }
    // Преобразуем variables в JSON-байты.
    jsonBytes, err := json.Marshal(ex.Variables)
    if err != nil {
        return err
    }

    jsonBytes2, err := json.Marshal(hist)
    if err != nil {
        return err
    }

    // Обновляем запись.
    query := `
        UPDATE example
        SET infix = ?, variables = ?, history = ?
        WHERE id = ?;`
    if _, err := db.ExecContext(ctx, query, infix, jsonBytes, jsonBytes2, id); err != nil {
        return err
    }

    return nil
}

func GetUserID(username string, id_ex int) (bool, error) {
    ctx := context.TODO()
    db, err := sql.Open("sqlite3", "database/calculator.db")
    if err != nil {
        log.Println("open db error")
        panic(err)
    }
    defer db.Close()

    err = db.PingContext(ctx)
    if err != nil {
        log.Println("connect error")
        panic(err)
    }
    row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)

    var id int
    err = row.Scan(&id)
    if err != nil {
        return false, err
    }

    idUser, err := GetUserIDFromExample(id_ex)
    if err != nil {
        return false, err
    }
    if idUser == id {
        return true, nil
    } else {
        return false, nil
    }
}


func GetUserIDFromExample(id int) (int, error) {
    ctx := context.TODO()
    db, err := sql.Open("sqlite3", "database/calculator.db")
    if err != nil {
        log.Println("open db error")
        panic(err)
    }
    defer db.Close()

    err = db.PingContext(ctx)
    if err != nil {
        log.Println("connect error")
        panic(err)
    }

    row := db.QueryRow("SELECT id_user FROM example WHERE id = ?", id)

    var idUser int
    err = row.Scan(&idUser)
    if err != nil {
        return 0, err
    }

    return idUser, nil
}

func GetHistoryAsJSON(id int) (Hist, error) {
    db, err := sql.Open("sqlite3", "database/calculator.db")
    if err != nil {
        return Hist{}, err
    }
    defer db.Close()

    row := db.QueryRow("SELECT history FROM example WHERE id = ?", id)

    var historyJSON string
    err = row.Scan(&historyJSON)
    if err != nil {
        return Hist{}, err
    }

    var history Hist
    err = json.Unmarshal([]byte(historyJSON), &history)
    if err != nil {
        return Hist{}, err
    }

    return history, nil
}
