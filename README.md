# go-api-test

Just a simple demonstration of creating an API in Go without any dependencies.

The HTTP method and route handling could be improved by using a package like `gin` to split routes by methods (`GET`, `POST`, etc) and to make handling named paths easier like `/api/users/:id`.

But the purpose of this project is simply to show what can be done with the standard `net/http` package.

## Concepts covered
- Read and write data to disk
- Fetch data from an API
- Encode and Decode JSON data
- Creating an http server and handling GET and POST requests
- Error handling and model validation

## Usage
Run `go run .` and navigate to `localhost:8080/api/users` in your browser to start

## Available Routes

- GET `/api/users`
- GET `/api/users/:id`
- GET `/api/todos`
- GET `/api/todos/:id`
- POST `/api/users`
- POST `/api/todos`