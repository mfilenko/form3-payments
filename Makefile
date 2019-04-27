.PHONY: dep test compose decompose start stop sandwich

compose_file = deployments/docker-compose.yml

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

compose:
	docker-compose -f $(compose_file) up --build --detach

decompose:
	docker-compose -f $(compose_file) down --volumes

start:
	docker-compose -f $(compose_file) start

stop:
	docker-compose -f $(compose_file) stop

sandwich:
	@[ $$(whoami) = "root" ] && (echo "ok"; echo "ham" > ~/sandwich) || (echo "make it yourself" && exit 13)
