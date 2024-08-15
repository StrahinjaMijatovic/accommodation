## syntax=docker/dockerfile:1
#
#FROM golang:1.20-alpine
#
#WORKDIR /app
#
#COPY go.mod ./
#COPY go.sum ./
#RUN go mod download
#
#COPY *.go ./
#
#RUN go build -o /accommodation-service
#
## Ensure the binary is executable
#RUN chmod +x /accommodation-service
#
#EXPOSE 8000
#
#CMD ["/accommodation-service"]


# syntax=docker/dockerfile:1
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /accommodation-service

EXPOSE 8080

CMD ["/accommodation-service"]
