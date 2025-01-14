FROM golang:latest as wallet
RUN apt-get update && apt-get install -y curl 
WORKDIR /app

# COPY ./pkg/proto/exchange ./pkg/proto
COPY ./gw-cyrrency-wallet ./
COPY ./pkg/config/config.yml ./config
COPY wallet.env ./

RUN go mod download
RUN sh -s -- -b $(go env GOPATH)/bin
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main ./cmd

CMD ["./main"]



FROM golang:latest as autstore
RUN apt-get update && apt-get install -y curl 
WORKDIR /app

COPY ./gw-exchanger ./
# COPY ./pkg/proto/exchange ./pkg/proto

RUN go mod download
RUN sh -s -- -b $(go env GOPATH)/bin
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main ./cmd

CMD ["./main"]