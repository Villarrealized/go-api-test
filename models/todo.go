package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todosCache []Todo

func FetchTodos() ([]Todo, error) {
	if todosCache != nil {
		fmt.Println("returning todos from cache...")
		return todosCache, nil
	}

	const filename string = "todos.json"

	todos, err := fetchTodosFromDisk(filename)
	if err != nil {
		todos, err = fetchTodosFromNetwork()
		if err != nil {
			return nil, err
		}

		err = saveData(todosCache, filename)
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

func fetchTodosFromDisk(filename string) ([]Todo, error) {
	data, err := readData(filename)
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
