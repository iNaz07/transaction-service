# syntax=docker/dockerfile:1
FROM golang:1.17-alpine as build

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN CGO_ENABLED=0 go build -o main ./app/main.go 

FROM alpine:latest

WORKDIR /cmd

COPY --from=build /app/main ./
COPY --from=build /app/config.json ./
COPY --from=build /app/templates ./templates

EXPOSE 8181

CMD ["./main"]
