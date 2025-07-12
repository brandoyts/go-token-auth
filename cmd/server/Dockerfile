FROM golang:1.23-alpine

RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

WORKDIR /app/cmd/server

ENTRYPOINT ["CompileDaemon", "--build=go build -o ./tmp/server", "--command=./tmp/server"]