package main

import (
	"encoding/json"
	"io"
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

func TestGetTodos(t *testing.T) {
	t.Run("it returns the todos as JSON", func(t *testing.T) {
		wantedTodos := []Todo{
			{"1", "Example 1", false},
			{"2", "Example 2", false},
			{"3", "Example 3", false},
		}

		store := StubTodoStore{wantedTodos}
		server := NewTodoServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/todos", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getTodosFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertTodos(t, got, wantedTodos)
	})
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
