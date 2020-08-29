package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubTodoStore struct {
	todos []Todo
}

func (s *StubTodoStore) GetTodos() []Todo {
	return s.todos
}

func (s *StubTodoStore) AddTodo(todo Todo) {
	s.todos = append(s.todos, todo)
}

func TestGetTodos(t *testing.T) {
	todos := []Todo{
		{"1", "Example 1", false},
		{"2", "Example 2", false},
		{"3", "Example 3", false},
	}
	store := StubTodoStore{todos}
	server := NewTodoServer(&store)

	t.Run("it returns the todos as JSON", func(t *testing.T) {
		request := newGetTodosRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getTodosFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertTodos(t, got, todos)
	})

	t.Run("it returns 400 if ID is missing", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/todos/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("it returns 404 on wrong ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/todos/4", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("it returns a todo as JSON", func(t *testing.T) {
		wantedTodo := Todo{"1", "Example 1", false}

		request, _ := http.NewRequest(http.MethodGet, "/todos/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getTodoFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertTodo(t, got, wantedTodo)
	})
}

func newGetTodosRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/todos", nil)
	return req
}

func newPostTodoRequest(todo Todo) *http.Request {
	json, err := json.Marshal(todo)
	if err != nil {
		log.Fatalf("could not marshal Todo into JSON. '%v'", err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewReader(json))
	return req
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func getTodoFromResponse(t *testing.T, body io.Reader) (todo Todo) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&todo)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into Todo, '%v'", body, err)
	}

	return
}

func getTodosFromResponse(t *testing.T, body io.Reader) (todos []Todo) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&todos)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Todo, '%v'", body, err)
	}

	return
}

func assertTodos(t *testing.T, got, want []Todo) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertTodo(t *testing.T, got, want Todo) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
