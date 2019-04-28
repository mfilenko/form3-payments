package main

func (server *Server) Setup() {
	server.Router.Handle("GET", "/payments", server.ReadPayments)
	server.Router.Handle("GET", "/payments/:id", server.ReadPayment)
	server.Router.Handle("POST", "/payments", server.CreatePayment)
	server.Router.Handle("PUT", "/payments/:id", server.UpgradePayment)
	server.Router.Handle("PATCH", "/payments/:id", server.UpdatePayment)
	server.Router.Handle("DELETE", "/payments/:id", server.DeletePayment)
}
