package servant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigStruct(t *testing.T) {
	host := "0.0.0.0"
	port := 8080
	loglevel := "Info"
	certPath := "foo.crt"
	keyPath := "foo.key"

	config := Config{
		Host:     host,
		Port:     port,
		Loglevel: loglevel,
		CertPath: certPath,
		KeyPath:  keyPath,
	}

	assert.Equal(t, config.Host, host)
	assert.Equal(t, config.Port, port)
	assert.Equal(t, config.Loglevel, loglevel)
	assert.Equal(t, config.CertPath, certPath)
	assert.Equal(t, config.KeyPath, keyPath)
}

func TestNewServant(t *testing.T) {
	config := Config{
		Host:     "0.0.0.0",
		Port:     8080,
		Loglevel: "Debug",
		CertPath: "foo.crt",
		KeyPath:  "foo.key",
	}


	auth := Auth{
		Enabled: true,
		Users: []User {
			User{Name: "Igor", PasswordSha512: "3fj39fd"},
		}}

route := NewRoute(
		"/api/foo",
		"GET",
		auth,
		func() string {return "WORKS!"},
)
	routes := []Route{route}

	_ = NewServant(config, routes)


	}
