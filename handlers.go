package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GET /payments - read all payments.
func ReadPayments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print(r.Method, r.URL.Path)
	fmt.Fprint(w, "Welcome!\n")
}

// GET /payments/:id - read a specific payment.
func ReadPayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print(r.Method, r.URL.Path)
	fmt.Fprint(w, "Welcome!\n")
}

// POST /payments - create new payment.
func CreatePayment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print(r.Method, r.URL.Path)
	fmt.Fprint(w, "Welcome!\n")
}

// PUT /payments/:id - replace payment completely.
func UpgradePayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print(r.Method, r.URL.Path)
	fmt.Fprint(w, "Welcome!\n")
}

// PATCH /payments/:id - update payment partially.
func UpdatePayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print(r.Method, r.URL.Path)
	fmt.Fprint(w, "Welcome!\n")
}

// DELETE /payments/:id - delete payment by ID.
func DeletePayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print(r.Method, r.URL.Path)
	fmt.Fprint(w, "Welcome!\n")
}
