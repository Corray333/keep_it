FROM golang:1.23-alpine3.20

WORKDIR /app

COPY . .
RUN apk add bash make musl-dev  
RUN go mod download
