FROM golang:1.23.3 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN go build -o main

FROM ubuntu:22.04
WORKDIR /app

COPY --from=build /app/main .
RUN apt-get update && apt-get install -y --no-install-recommends default-mysql-client && apt-get clean
EXPOSE 8080
CMD ["./main"]
