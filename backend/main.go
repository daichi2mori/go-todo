package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
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
		content TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todo", getTodo)
	mux.HandleFunc("POST /todo", createTodo)
	mux.HandleFunc("PUT /todo", updateTodo)
	mux.HandleFunc("DELETE /todo/{id}", deleteTodo)

	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	cmd := `SELECT * FROM todos`

	rows, err := db.Query(cmd)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Content, &todo.Completed)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todo.Content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	cmd := `INSERT INTO todos (content, completed) VALUES (?, ?)`
	_, err = db.Exec(cmd, todo.Content, todo.Completed)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create Todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo created successfully"})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todo.ID == 0 {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	if todo.Content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	cmd := `UPDATE todos SET content = ?, completed = ? WHERE id = ?`
	_, err = db.Exec(cmd, todo.Content, todo.Completed, todo.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update Todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo updated successfully"})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	// URLパスからIDを取得
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// IDがデータベースに存在するか確認
	var exists bool
	cmd := "SELECT EXISTS(SELECT 1 FROM todos WHERE id = ?)"
	err = db.QueryRow(cmd, id).Scan(&exists)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to check existence of Todo", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// データベースから削除
	cmd = `DELETE FROM todos WHERE id = ?`
	_, err = db.Exec(cmd, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete Todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
}
