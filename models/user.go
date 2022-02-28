package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}

var usersCache []User

func FetchUsers() ([]User, error) {
	if usersCache != nil {
		fmt.Println("returning users from cache...")
		return usersCache, nil
	}

	var filename string = "users.json"

	usersCache, err := fetchUsersFromDisk(filename)
	if err != nil {
		usersCache, err = fetchUsersFromNetwork()
		if err != nil {
			return nil, err
		}

		err = saveData(usersCache, filename)
		if err != nil {
			return nil, err
		}
	}

	return usersCache, nil
}

func fetchUsersFromDisk(filename string) ([]User, error) {
	data, err := readData(filename)
	if err != nil {
		return nil, err
	}

	var users []User

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	fmt.Println("returning users from disk...")
	return users, nil
}

func fetchUsersFromNetwork() ([]User, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/users")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []User

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&users)

	if err != nil {
		return nil, err
	}

	fmt.Println("returning users from network...")
	return users, nil
}
