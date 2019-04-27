.PHONY: dep test

default:
	go build

dep:
	go get -u github.com/go-openapi/strfmt
	go get -u github.com/julienschmidt/httprouter
	go get -u github.com/lib/pq
	go get -u github.com/jmoiron/sqlx
	go get -u github.com/kelseyhightower/envconfig


test:
	go test -v
