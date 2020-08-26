package main

import (
	"encoding/json"
	"log"
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
	AddTodo(todo Todo)
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
	switch r.Method {
	case http.MethodPost:
		var todo Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			log.Fatalf("Unable to parse request body %q, '%v'", r.Body, err)
		}
		t.store.AddTodo(todo)
		w.WriteHeader(http.StatusAccepted)
	case http.MethodGet:
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(t.store.GetTodos())
	}
}
