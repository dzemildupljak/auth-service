FROM golang:alpine AS builder

RUN apk update && apk add openssl

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN touch log/debug-log.log
RUN touch log/error-log.log


RUN openssl genrsa -out access-private.pem 2048
RUN openssl rsa -in access-private.pem -outform PEM -pubout -out access-public.pem

RUN openssl genrsa -out refresh-private.pem 2048
RUN openssl rsa -in refresh-private.pem -outform PEM -pubout -out refresh-public.pem

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service

# Run stage
FROM alpine

COPY --from=builder /app/auth-service /app/

COPY --from=builder /app/access-private.pem ./
COPY --from=builder /app/refresh-private.pem ./

COPY --from=builder /app/access-public.pem ./
COPY --from=builder /app/refresh-public.pem ./

RUN mkdir -p /log


EXPOSE 80

ENTRYPOINT ["/app/auth-service"]