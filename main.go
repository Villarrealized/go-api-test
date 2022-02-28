package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/villarrealized/go-api-test/models"
)

const usersUrl string = "/api/users/"
const todosUrl string = "/api/todos/"

type jsonError struct {
	Message string `json:"error"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(usersUrl, handleUsers)
	mux.HandleFunc(todosUrl, handleTodos)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func handleUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	stringId := strings.TrimPrefix(request.URL.Path, usersUrl)

	id, err := strconv.Atoi(stringId)
	if err == nil {
		writeUser(id, writer, request)
		return
	}

	users, err := models.FetchUsers()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(writer, users)
}

func handleTodos(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	stringId := strings.TrimPrefix(request.URL.Path, todosUrl)

	id, err := strconv.Atoi(stringId)
	if err == nil {
		writeTodo(id, writer, request)
		return
	}

	todos, err := models.FetchTodos()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(writer, todos)
}

func writeUser(id int, writer http.ResponseWriter, request *http.Request) {
	user, err := models.FetchUser(id)
	if err != nil {
		handleFetchRecordError(err, writer)
		return
	}

	writeJson(writer, user)
}

func writeTodo(id int, writer http.ResponseWriter, request *http.Request) {
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
