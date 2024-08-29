package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

var db *sql.DB

// func init() {
// 	var err error
// 	db, err = sql.Open("sqlite3", "todos.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	createTable := `
// 	CREATE TABLE IF NOT EXISTS todos (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		content TEXT,
// 		Completed BOOLEAN
// 	);`

// 	_, err = db.Exec(createTable)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todo", test)
	mux.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "POST"})
	})
	mux.HandleFunc("PUT /todo", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "PUT"})
	})
	mux.HandleFunc("DELETE /todo", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "DELETE"})
	})

	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "GET"})
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	cmd := `SELECT * FROM todos`

	rows, err := db.Query(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Content, &todo.Done)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	defer rows.Close()

	// wのヘッダーにcontent-typeを追加
	w.Header().Set("Content-Type", "application/json")
	// NewEncoder()の引数内にエンコーダーを作成
	// Encode()の引数を実際にエンコードする
	json.NewEncoder(w).Encode(todos)
}
