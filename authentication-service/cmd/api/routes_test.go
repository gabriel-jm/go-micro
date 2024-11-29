package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_routes_exist(t *testing.T) {
	testApp := Config{}
	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	routes := []string{"/authenticate"}

	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, router chi.Router, route string) {
	found := false

	chi.Walk(
		router,
		func(method, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			if route == foundRoute {
				found = true
			}

			return nil
		},
	)

	if !found {
		t.Errorf("Did not find %s in registered routes", route)
	}
}