package main

// NewInMemoryTodoStore initializes an empty todo store.
func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{[]Todo{}}
}

// InMemoryTodoStore holds todos in memory.
type InMemoryTodoStore struct {
	store []Todo
}

// GetTodos returns Todos from memory.
func (i *InMemoryTodoStore) GetTodos() []Todo {
	return i.store
}

// AddTodo adds Todo to memory.
func (i *InMemoryTodoStore) AddTodo(todo Todo) {
	i.store = append(i.store, todo)
}
