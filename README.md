# go-api-test

Just a simple demonstration of creating an API in Go without any dependencies.

The HTTP method and route handling could be improved by using a package like `gin` to split routes by methods (`GET`, `POST`, etc) and to make handling named paths easier like `/api/users/:id`.

But the purpose of this project is simply to show what can be done with the standard `net/http` package. 