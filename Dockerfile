FROM golang:1.18 as development


WORKDIR /app

COPY . /app

RUN touch log/debug-log.log
RUN touch log/error-log.log


RUN go mod download

RUN go build -o auth-service-img main.go

EXPOSE 8004

# CMD [ "./auth-service" ]
