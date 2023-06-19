FROM --platform=linux/amd64 golang:1.20

WORKDIR /app

COPY go.mod go.sum . 

RUN go mod download

COPY . .

RUN go build -o ./bin/spotify-party

EXPOSE 8080

CMD ["/app/bin/spotify-party"]
