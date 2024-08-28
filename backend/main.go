package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Todo struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "todos.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		Completed BOOLEAN
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

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
	//
	json.NewEncoder(w).Encode(todos)
}
