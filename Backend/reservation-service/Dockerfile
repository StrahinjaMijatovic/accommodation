FROM golang:1.20-alpine

WORKDIR /app

# Instaliraj bash
RUN apk add --no-cache bash

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY wait-for-it.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/wait-for-it.sh

RUN go build -o /reservation-service

EXPOSE 8080

CMD ["wait-for-it.sh", "neo4j:7687", "--", "/reservation-service"]
