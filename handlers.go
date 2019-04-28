package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/go-openapi/strfmt"
	"github.com/imdario/mergo"
)

type Response struct {
	Data  interface{} `json:"data"`
	Links `json:"links"`
}

type Links struct {
	Self string `json:"self"`
}

type Error struct {
	Message string `json:"message"`
}

// Response boilerplate code.
func serve(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if code == 204 || code >= 400 {
		json.NewEncoder(w).Encode(&Error{Message: http.StatusText(code)})
	} else {
		json.NewEncoder(w).Encode(data)
	}
}

// GET /payments - read all payments.
func (server *Server) ReadPayments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var status int
	var response *Response

	if payments, err := GetPayments(server.DB); err != nil {
		log.Print(err)
		status = http.StatusInternalServerError
	} else {
		status = http.StatusOK
		response = &Response{
			Data:  payments,
			Links: Links{Self: "https://api.test.form3.tech/v1/payments"},
		}
	}

	serve(w, status, response)
}

// GET /payments/:id - read a specific payment.
func (server *Server) ReadPayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var status int
	var response *Response

	payment := Payment{ID: strfmt.UUID4(ps.ByName("id"))}

	if err := payment.Get(server.DB); err != nil {
		log.Print(err)
		switch err {
		case sql.ErrNoRows:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
	} else {
		status = http.StatusOK
		response = &Response{
			Data:  payment,
			Links: Links{Self: "https://api.test.form3.tech/v1/payments"},
		}
	}

	serve(w, status, response)
}

// POST /payments - create new payment.
func (server *Server) CreatePayment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var status int
	var response *Response
	var payment Payment

	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		log.Print(err)
		status = http.StatusBadRequest
	} else if err := payment.Save(server.DB); err != nil {
		log.Print(err)
		status = http.StatusInternalServerError
	} else {
		status = http.StatusCreated
		response = &Response{
			Data:  payment,
			Links: Links{Self: "https://api.test.form3.tech/v1/payments"},
		}
	}

	serve(w, status, response)
}

// PUT /payments/:id - replace payment completely.
func (server *Server) UpgradePayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var status int
	var response *Response
	var payload Payment

	payment := Payment{ID: strfmt.UUID4(ps.ByName("id"))}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Print(err)
		status = http.StatusBadRequest
	} else if err := payment.Get(server.DB); err != nil {
		log.Print(err)
		switch err {
		case sql.ErrNoRows:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
	} else {
		payload.ID = payment.ID
		payload.Version = payment.Version + 1
		if err := payload.Overwrite(server.DB); err != nil {
			log.Print(err)
			status = http.StatusInternalServerError
		} else {
			status = http.StatusOK
			response = &Response{
				Data:  payload,
				Links: Links{Self: "https://api.test.form3.tech/v1/payments"},
			}
		}
	}

	serve(w, status, response)
}

// PATCH /payments/:id - update payment partially.
func (server *Server) UpdatePayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var status int
	var response *Response
	var payload Payment

	payment := Payment{ID: strfmt.UUID4(ps.ByName("id"))}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Print(err)
		status = http.StatusBadRequest
	} else if err := payment.Get(server.DB); err != nil {
		log.Print(err)
		switch err {
		case sql.ErrNoRows:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
	} else {
		mergo.Merge(&payload, payment)
		payload.Version++
		if err := payload.Overwrite(server.DB); err != nil {
			log.Print(err)
			status = http.StatusInternalServerError
		} else {
			status = http.StatusOK
			response = &Response{
				Data:  payload,
				Links: Links{Self: "https://api.test.form3.tech/v1/payments"},
			}
		}
	}

	serve(w, status, response)
}

// DELETE /payments/:id - delete payment by ID.
func (server *Server) DeletePayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var status int

	payment := Payment{ID: strfmt.UUID4(ps.ByName("id"))}

	if err := payment.Remove(server.DB); err != nil {
		log.Print(err)
		status = http.StatusInternalServerError
	} else {
		status = http.StatusNoContent
	}

	serve(w, status, nil)
}
