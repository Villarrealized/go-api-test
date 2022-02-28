package main

import (
	"fmt"

	"github.com/villarrealized/go-api-test/models"
)

func main() {
	_, err := models.FetchUsers()

	if err != nil {
		fmt.Println(err)
	}

	_, err = models.FetchTodos()

	if err != nil {
		fmt.Println(err)
	}
}
