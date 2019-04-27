package main

import "github.com/julienschmidt/httprouter"

func Router(routes Routes) *httprouter.Router {
	router := httprouter.New()

	for _, route := range routes {
		var handle httprouter.Handle
		handle = Logger(route.Handle)
		router.Handle(route.Method, route.Path, handle)
	}

	return router
}
