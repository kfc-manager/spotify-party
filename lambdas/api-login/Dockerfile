FROM golang:1.21.2-alpine3.18 as build

WORKDIR /app

COPY go.mod go.sum ./

COPY main.go ./

RUN go build -o main main.go

FROM alpine:3.18

COPY --from=build /app/main /main

ENTRYPOINT [ "/main" ]
