package main

import (
	"encoding/json"
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
		Addr:    ":8080",
		Handler: mux,
	}
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
			return
		}

		writeJson(writer, users)
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
			return
		}

		writeJson(writer, todos)
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
			return
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			writeJson(writer, jsonError{"Internal server error"})
			return
		}
	}
}

func writeJson(writer http.ResponseWriter, data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonBytes)
}
