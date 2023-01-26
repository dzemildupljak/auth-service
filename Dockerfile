FROM golang:1.18 as development


WORKDIR /app

COPY . /app

RUN go mod download

# Install Reflex for development
# RUN go install github.com/cespare/reflex@latest
RUN go build main.go .

EXPOSE 8004


CMD [ "./auth-service" ]
# CMD reflex -r '\.go$$' -s go run *.go