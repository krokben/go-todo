package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingTodosAndRetrievingThem(t *testing.T) {
	store := NewInMemoryTodoStore()
	server := NewTodoServer(store)

	todo := Todo{"1", "Example task", false}

	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest(todo))
	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest(todo))
	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest(todo))

	t.Run("get todos", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetTodosRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getTodosFromResponse(t, response.Body)
		want := []Todo{todo, todo, todo}
		assertTodos(t, got, want)
	})
}
