package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello interweb")
}

func what(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "sup")
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gotest")
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
