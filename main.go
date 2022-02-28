package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/villarrealized/go-api-test/models"
)

const (
	usersUrl string = "/api/users"
	// For users/:id
	userUrl string = "/api/users/"
)

const (
	todosUrl string = "/api/todos"
	// For todos/:id
	todoUrl string = "/api/todos/"
)

type jsonError struct {
	Message string `json:"error"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(usersUrl, UsersHandler)
	mux.HandleFunc(userUrl, UserHandler)

	mux.HandleFunc(todosUrl, TodosHandler)
	mux.HandleFunc(todoUrl, TodoHandler)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	fmt.Printf("Listening on %s\n", server.Addr)
	server.ListenAndServe()
}

func UsersHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	switch request.Method {
	case http.MethodGet:
		users, err := models.FetchUsers()
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			writeJson(writer, jsonError{"Internal server error"})
			return
		}

		writeJson(writer, users)
		return

	case http.MethodPost:
		var user models.User
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&user)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			writeJson(writer, jsonError{"Invalid json data"})
			return
		}

		newUser, err := models.CreateUser(user)
		if err != nil {
			log.Println(err)
			switch typedErr := err.(type) {
			case *models.ModelMissingRequiredFieldError:
				writer.WriteHeader(http.StatusUnprocessableEntity)
				writeJson(writer, jsonError{typedErr.Message})
			case *models.UniqueViolationError:
				writer.WriteHeader(http.StatusUnprocessableEntity)
				writeJson(writer, jsonError{typedErr.Message})
			default:
				writer.WriteHeader(http.StatusInternalServerError)
				writeJson(writer, jsonError{"Internal server error"})
			}
			return
		}

		writeJson(writer, newUser)
		return

	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(writer, jsonError{"Method not allowed"})
	}
}

// Handling for api/users/:id
func UserHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	stringId := strings.TrimPrefix(request.URL.Path, userUrl)

	id, err := strconv.Atoi(stringId)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writeJson(writer, jsonError{"Invalid id"})
		return
	}

	switch request.Method {
	case http.MethodGet:
		writeUser(writer, id)
		return
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(writer, jsonError{"Method not allowed"})
		return
	}
}

func TodosHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	switch request.Method {
	case http.MethodGet:
		todos, err := models.FetchTodos()
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			writeJson(writer, jsonError{"Internal server error"})
			return
		}

		writeJson(writer, todos)
		return

	case http.MethodPost:
		var todo models.Todo
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&todo)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			writeJson(writer, jsonError{"Invalid json data"})
			return
		}

		newTodo, err := models.CreateTodo(todo)
		if err != nil {
			log.Println(err)
			switch typedErr := err.(type) {
			case *models.ModelMissingRequiredFieldError:
				writer.WriteHeader(http.StatusUnprocessableEntity)
				writeJson(writer, jsonError{typedErr.Message})
			case *models.ModelRelationshipError:
				writer.WriteHeader(http.StatusUnprocessableEntity)
				writeJson(writer, jsonError{typedErr.Message})
			default:
				writer.WriteHeader(http.StatusInternalServerError)
				writeJson(writer, jsonError{"Internal server error"})
			}
			return
		}

		writeJson(writer, newTodo)
		return

	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(writer, jsonError{"Method not allowed"})
	}
}

// Handling for api/todos/:id
func TodoHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	stringId := strings.TrimPrefix(request.URL.Path, todoUrl)

	id, err := strconv.Atoi(stringId)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writeJson(writer, jsonError{"Invalid id"})
		return
	}

	switch request.Method {
	case http.MethodGet:
		writeTodo(writer, id)
		return
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(writer, jsonError{"Method not allowed"})
		return
	}
}

func writeUser(writer http.ResponseWriter, id int) {
	user, err := models.FetchUser(id)
	if err != nil {
		handleFetchRecordError(err, writer)
		return
	}

	writeJson(writer, user)
}

func writeTodo(writer http.ResponseWriter, id int) {
	todo, err := models.FetchTodo(id)
	if err != nil {
		handleFetchRecordError(err, writer)
		return
	}

	writeJson(writer, todo)
}

func handleFetchRecordError(err error, writer http.ResponseWriter) {
	if err != nil {
		log.Println(err)
		switch typedErr := err.(type) {
		case *models.ModelMissingError:
			writer.WriteHeader(http.StatusNotFound)
			writeJson(writer, jsonError{typedErr.Message})
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			writeJson(writer, jsonError{"Internal server error"})
		}
	}
}

func writeJson(writer http.ResponseWriter, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "error": "Internal server error" }`))
		return
	}

	writer.Write(jsonBytes)
}
