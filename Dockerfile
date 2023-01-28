FROM golang:1.18 as development


WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o auth-service main.go

EXPOSE 8004
