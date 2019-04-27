package main

import "github.com/julienschmidt/httprouter"

func Router(routes Routes) *httprouter.Router {
	router := httprouter.New()

	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handle)
	}

	return router
}
