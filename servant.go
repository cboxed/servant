package servant

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cboxed/servant/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Config holds the general config for servant
type Config struct {
	Host     string
	Port     int
	Loglevel string
	TLS      bool
	CertPath string
	KeyPath  string
}

// Servant is the servant object
type Servant struct {
	config Config
	routes []Route
}

// NewServant creates a new servant instance
func NewServant(config Config, routes []Route) Servant {
	return Servant{config: config, routes: routes}
}

func (s *Servant) prepareCerts() {
	if s.config.CertPath == "" {
		s.config.CertPath = "server.pem"
	}

	if s.config.KeyPath == "" {
		s.config.KeyPath = "server.key"
	}

	utils.CreateCertificatesIfTheyDoNotExist(s.config.CertPath, s.config.KeyPath)
}

// Summon starts the servant instance to listen
func (s *Servant) Summon() {

	router := mux.NewRouter()

	// setup routes based on routes config
	log.Info("add routes")
	for _, route := range s.routes {
		router.HandleFunc(route.Endpoint, s.authHandler(route.Serve)).Methods(route.Method)
		log.Info("  --> " + route.Method + " " + route.Endpoint)
	}

	// must be at the end. catches everything what not matched before
	router.PathPrefix("/").HandlerFunc(s.authHandler(s.defaultHandler))
	log.Info("  --> DefaultRoute, catches everything what not yet matched.")

	hostPort := s.config.Host + ":" + strconv.Itoa(s.config.Port)

	log.Info("listen on " + hostPort)

	if s.config.TLS {
		// TLS is enabled
		log.Info("TLS enabled")
		s.prepareCerts()
		log.Fatal(http.ListenAndServeTLS(hostPort, s.config.CertPath, s.config.KeyPath, handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)))
	} else {
		// TLS is disabled
		log.Info("TLS disabled")
		log.Fatal(http.ListenAndServe(hostPort, handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)))
	}
}

func (s *Servant) defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(fmt.Sprintf("unknown endpoint called: %s (%s)", r.URL, r.RemoteAddr))
	http.Error(w, "Endpoint does not exist", 404)
}

func (s *Servant) authHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// authentication is based on called route
		url := fmt.Sprintf("%#v", r.URL)
		routeConfig := s.getEndpointConfig(url)
		if !routeConfig.Auth.Enabled {
			log.Info(fmt.Sprintf("authentication is disabled for endpoint: %s (%s)", r.URL, r.RemoteAddr))
			f(w, r)
			return
		}

		log.Info(fmt.Sprintf("authHandler called for endpoint: %s (%s)", r.URL, r.RemoteAddr))

		// get fields from authorization header
		// assume fields[0]=user and fields[1]=sha512(password)
		fields := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(fields) != 2 {
			log.Info(fmt.Sprintf("malformed authentication header received (%s)", r.RemoteAddr))
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(fields[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		pair := strings.SplitN(string(b), ":", 2)

		// not exactly two parameters were send username:password
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		// password stored in local config is sha512 of real password
		// real password is transmitted and we generate the sha512 of it
		if len(pair) == 2 {
			log.Info(fmt.Sprintf("try to authenticate user: %s (%s)", pair[0], r.RemoteAddr))
		}

		// check if user is allowed
		if s.userIsAllowed(pair[0], pair[1], routeConfig) {
			log.Info(fmt.Sprintf("authentication successful for user: %s (%s)", pair[0], r.RemoteAddr))
			f(w, r)
		} else {
			log.Info(fmt.Sprintf("authentication failed for user: %s (%s)", pair[0], r.RemoteAddr))
			http.Error(w, "Not authorized", 401)
			return
		}
	}
}

func (s *Servant) userIsAllowed(username string, password string, routeConfig Route) bool {

	// if authentication is disabled for this route then everybody is allowed
	if !routeConfig.Auth.Enabled {
		return true
	}

	passwordHash := strings.ToLower(utils.StringToSha512(password))

	// find the correct user
	for _, user := range routeConfig.Auth.Users {
		if user.Name == username {
			// check password
			if passwordHash == strings.ToLower(user.PasswordSha512) {
				return true
			}
		}
	}
	// user is not in list of allowed users for this route
	return false
}

//func (s Servant) getEndpointConfig(endpoint string) Route {
func (s *Servant) getEndpointConfig(endpoint string) Route {
	for _, route := range s.routes {
		if endpoint == route.Endpoint {
			return route
		}
	}

	// return default route
	return Route{
		Endpoint: "/",
		Method:   "GET",
		Auth: Auth{
			Enabled: false,
			Users:   []User{},
		},
		Function: func() string { return "bla" },
	}
}
