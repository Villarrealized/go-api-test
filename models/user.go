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

/*
Tries to return data first from the in-memory
cache, then disk, and then the network.
*/
func FetchUsers() ([]User, error) {
	if usersCache != nil {
		fmt.Println("returning users from cache...")
		return usersCache, nil
	}

	const filename string = "users.json"

	users, err := fetchUsersFromDisk(filename)
	if err != nil {
		users, err = fetchUsersFromNetwork()
		if err != nil {
			return nil, err
		}

		err = saveData(usersCache, filename)
		if err != nil {
			return nil, err
		}
	}

	usersCache = users

	return usersCache, nil
}

func FetchUser(id int) (*User, error) {
	users, err := FetchUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Id == id {
			return &user, nil
		}
	}

	return nil, &ModelMissingError{"No user found for that id"}
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
