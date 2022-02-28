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

const usersFile string = "users.json"

/*
Tries to return data first from the in-memory
cache, then disk, and then the network.
*/
func FetchUsers() ([]User, error) {
	if usersCache != nil {
		fmt.Println("returning users from cache...")
		return usersCache, nil
	}

	users, err := fetchUsersFromDisk(usersFile)
	if err != nil {
		users, err = fetchUsersFromNetwork()
		if err != nil {
			return nil, err
		}

		err = saveData(usersCache, usersFile)
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

func CreateUser(newUser User) (*User, error) {
	users, err := FetchUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Username == newUser.Username {
			return nil, &UniqueViolationError{"That username already exists."}
		}
	}

	id, err := getNextUserId()
	if err != nil {
		return nil, err
	}
	newUser.Id = id

	// Update the cache
	usersCache = append(usersCache, newUser)

	// Save to disk
	err = saveData(usersCache, usersFile)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// Implementation is potentially not safe...could have conflicting ids
func getNextUserId() (int, error) {
	users, err := FetchUsers()
	if err != nil {
		return 0, err
	}

	return users[len(users)-1].Id + 1, nil
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
