package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/villarrealized/go-api-test/models"
)

var users []models.User

func main() {
	loadData()
}

func loadData() {

	users, err := models.FetchUsers()

	if err != nil {
		fmt.Println(err)
	}

	err = saveData(users)
	if err != nil {
		fmt.Println(err)
	}
}

func saveData(data interface{}) error {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	var filename string

	switch dataType := data.(type) {
	case []User:
		filename = "users.json"
	default:
		return fmt.Errorf("Unknown type: %v", dataType)
	}

	err = ioutil.WriteFile("data/"+filename, jsonData, 0664)
	if err != nil {
		return err
	}

	return nil
}
