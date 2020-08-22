package main

import (
	"encoding/json"
	"net/http"
)

// Todo stores a task that has an id and is either completed or not.
type Todo struct {
	ID        string
	Task      string
	Completed bool
}

// TodoStore stores Todos.
type TodoStore interface {
	GetTodos() []Todo
}

// TodoServer is an HTTP interface for Todos.
type TodoServer struct {
	store TodoStore
	http.Handler
}

const jsonContentType = "application/json"

// NewTodoServer creates a TodoServer with routing configured.
func NewTodoServer(store TodoStore) *TodoServer {
	t := new(TodoServer)

	t.store = store

	router := http.NewServeMux()
	router.Handle("/todos", http.HandlerFunc(t.todosHandler))

	t.Handler = router

	return t
}

func (t *TodoServer) todosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(t.store.GetTodos())
}
