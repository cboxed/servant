package servant

import (
	"fmt"
	"net/http"
)

type User struct {
	Name           string
	PasswordSha512 string
}

type Auth struct {
	Enabled bool
	Users   []User
}

type Route struct {
	Endpoint string
	Method   string
	Auth     Auth
	Function func() string
}

func NewRoute(endpoint string, method string, auth Auth, function func() string) Route {
	return Route{Endpoint: endpoint, Method: method, Function: function, Auth: auth}
}

func (route Route) Serve(w http.ResponseWriter, r *http.Request) {
	data := route.Function()
	_, _ = fmt.Fprintln(w, data)
}
