package servant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserStruct(t *testing.T) {


	name := "Igor"
	passHash := "3fj39fd"
	u:=User{
		Name:           name,
		PasswordSha512: passHash,
	}

	assert.Equal(t, u.Name, name)
	assert.Equal(t, u.PasswordSha512, passHash)
}

func TestAuthStruct(t *testing.T) {
	state := true
	user := User{
		Name:           "Igor",
		PasswordSha512: "3fj39fd",
	}
	users := []User {user}

	auth := Auth{
		Enabled: state,
		Users:   users,
	}

	assert.Equal(t, auth.Enabled, state)
	assert.Equal(t, auth.Users, []User {User{
		Name:           "Igor",
		PasswordSha512: "3fj39fd",
	}})
}

func TestRouteStruct(t *testing.T) {
	endpoint := "/foobar"
	method := "GET"
	user := User{
		Name:           "Igor",
		PasswordSha512: "3fj39fd",
	}
	users := []User {user}

	auth := Auth{
		Enabled: true,
		Users:   users,
	}

	route := Route{
		Endpoint: endpoint,
		Method:   method,
		Auth:     auth,
		Function: func() string {return "string"},
	}

	assert.Equal(t, route.Endpoint, endpoint)
	assert.Equal(t, route.Method, method)
	assert.Equal(t, route.Auth, auth)

}
