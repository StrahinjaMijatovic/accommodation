FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /accommodation-service

EXPOSE 8080

CMD ["/accommodation-service"]