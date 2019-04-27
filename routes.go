package main

import "github.com/julienschmidt/httprouter"

type Route struct {
	Method string
	Path   string
	Handle httprouter.Handle
}

type Routes []Route

func Endpoints() Routes {
	routes := Routes{
		Route{"GET", "/payments", ReadPayments},
		Route{"GET", "/payments/:id", ReadPayment},
		Route{"POST", "/payments", CreatePayment},
		Route{"PUT", "/payments/:id", UpgradePayment},
		Route{"PATCH", "/payments/:id", UpdatePayment},
		Route{"DELETE", "/payments/:id", DeletePayment},
	}
	return routes
}
