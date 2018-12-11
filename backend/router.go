package backend

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const (
	// APIVersion defines the compatability version of the API and is appended to each API route
	APIVersion     = "1"
	endpointFormat = "/api/v%s/%s"
)

// getEndpoint returns a properly formatted API endpoint
func getEndpoint(suffix string) string {
	return fmt.Sprintf(endpointFormat, APIVersion, suffix)
}

// Route defines a route passed to our mux
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes holds a list of Routes
type Routes []Route

// all defined server endpoints
var routes = Routes{

	// API endpoints
	Route{
		"Version",
		"GET",
		getEndpoint("version"),
		handlerVersion,
	},
	Route{
		"Refresh",
		"POST",
		getEndpoint("refresh"),
		handlerRefresh,
	},
	Route{
		"EnvAll",
		"GET",
		getEndpoint("env/{group:summary|details}"),
		handlerEnvAll,
	},
	Route{
		"EnvSingle",
		"GET",
		getEndpoint("env/{env-id}/{group:summary|details}"),
		handlerEnvSingle,
	},
	Route{
		"EnvPowerToggle",
		"POST",
		getEndpoint("env/{env-id}/{state:start|stop}"),
		handlerEnvPowerToggle,
	},
	Route{
		"InstancePowerToggle",
		"POST",
		getEndpoint("instance/{instance-id}/{state:start|stop}"),
		handlerInstancePowerToggle,
	},
}

func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc

		// add routes to mux
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// add route to mux to handle frontend UI static files (generated by npm)
	staticPath := viper.GetString("server.static_files_dir")
	if staticPath == "" {
		staticPath = "./frondent/dist"
	}

	router.
		Methods("GET").
		PathPrefix("/").
		Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticPath))))

	return router
}