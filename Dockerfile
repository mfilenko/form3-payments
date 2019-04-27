FROM golang:1 as builder
WORKDIR /go/src/github.com/mfilenko/form3/payments
RUN go get -d -v github.com/go-openapi/strfmt
RUN go get -d -v github.com/julienschmidt/httprouter
RUN go get -d -v github.com/lib/pq
RUN go get -d -v github.com/jmoiron/sqlx
RUN go get -d -v github.com/kelseyhightower/envconfig
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payments .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/mfilenko/form3/payments/payments .
CMD ["./payments"]
