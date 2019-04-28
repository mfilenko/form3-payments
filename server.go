package main

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Router *httprouter.Router
	DB     *sqlx.DB
	Config Configuration
}

func NewServer() *Server {
	return &Server{
		Config: Configure(),
	}
}

func (server *Server) Start() {
	server.Router = httprouter.New()
	server.Setup()
	server.OpenDB()

	log.Fatal(http.ListenAndServe(server.Config.Addr, server.Router))
}

func (server *Server) Stop() {
	server.DB.Close()
}
