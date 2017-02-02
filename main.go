package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello interweb")
}

func what(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "sup")
}

func main() {

	connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user := mustGetenv("CLOUDSQL_USER")
	password := os.Getenv("CLOUDSQL_PASSWORD")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName))
	if err != nil {
		fmt.Println("Something connect open")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("stuff didnt work")
	}
	fmt.Println("stuff")
	var name string
	var quantity int
	rows, err := db.Query("SELECT name, quantity FROM godata")
	if err != nil {
		fmt.Print("fail in select")
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &quantity)
		if err != nil {
			fmt.Println("fel i scan")
		}
		fmt.Println(name, quantity)
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/sup", what)
	http.ListenAndServe(":8000", nil)
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set ", k)
	}
	return v
}
