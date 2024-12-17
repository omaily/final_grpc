FROM golang:latest as wallet

RUN apt-get update && apt-get install -y curl 
WORKDIR /app

COPY ./ ./

RUN go mod download
RUN sh -s -- -b $(go env GOPATH)/bin
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main ./cmd

CMD ["./main"]