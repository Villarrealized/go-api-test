package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/villarrealized/go-api-test/helpers"
)

type Todo struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todosCache []Todo

const todosFile string = "todos.json"

/*
Tries to return data first from the in-memory
cache, then disk, and then the network.
*/
func FetchTodos() ([]Todo, error) {
	if todosCache != nil {
		fmt.Println("returning todos from cache...")
		return todosCache, nil
	}

	todos, err := fetchTodosFromDisk(todosFile)
	if err != nil {
		todos, err = fetchTodosFromNetwork()
		if err != nil {
			return nil, err
		}

		err = helpers.SaveData(todos, todosFile)
		if err != nil {
			return nil, err
		}
	}

	todosCache = todos

	return todosCache, nil
}

func FetchTodo(id int) (*Todo, error) {
	todos, err := FetchTodos()
	if err != nil {
		return nil, err
	}

	for _, todo := range todos {
		if todo.Id == id {
			return &todo, nil
		}
	}

	return nil, &ModelMissingError{"No todo found for that id"}
}

func CreateTodo(newTodo Todo) (*Todo, error) {

	if newTodo.UserId <= 0 {
		return nil, &ModelMissingRequiredFieldError{"userId field is required"}
	}
	if newTodo.Title == "" {
		return nil, &ModelMissingRequiredFieldError{"title field is required"}
	}

	users, err := FetchUsers()
	if err != nil {
		return nil, err
	}
	var foundUser bool = false
	for _, user := range users {
		if user.Id == newTodo.UserId {
			foundUser = true
		}
	}

	if !foundUser {
		return nil, &ModelRelationshipError{"That userId does not exist."}
	}

	id, err := getNextTodoId()
	if err != nil {
		return nil, err
	}
	newTodo.Id = id

	// Update the cache
	todosCache = append(todosCache, newTodo)

	// Save to disk
	err = helpers.SaveData(todosCache, todosFile)
	if err != nil {
		return nil, err
	}

	return &newTodo, nil
}

// Implementation is potentially not safe...could have conflicting ids
func getNextTodoId() (int, error) {
	todos, err := FetchTodos()
	if err != nil {
		return 0, err
	}

	return todos[len(todos)-1].Id + 1, nil
}

func fetchTodosFromDisk(filename string) ([]Todo, error) {
	data, err := helpers.ReadData(filename)
	if err != nil {
		return nil, err
	}

	var todos []Todo

	err = json.Unmarshal(data, &todos)
	if err != nil {
		return nil, err
	}

	fmt.Println("returning todos from disk...")
	return todos, nil
}

func fetchTodosFromNetwork() ([]Todo, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var todos []Todo

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&todos)

	if err != nil {
		return nil, err
	}

	fmt.Println("returning todos from network...")
	return todos, nil
}
