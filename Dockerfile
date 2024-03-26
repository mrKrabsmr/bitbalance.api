FROM golang:1.22.1

WORKDIR /api-server

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

EXPOSE 8000


