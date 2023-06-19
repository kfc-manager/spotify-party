package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kfc-manager/spotify-party/server"
)

func main() {
    godotenv.Load() // load environment variables

    port := os.Getenv("PORT")

    if len(port) < 1 { // incase the port is not defined in .env
        port = "8080" // we use the default port of 8080
    }

    s := server.New(":" + os.Getenv("PORT")) // create the server

    err := s.Run() // initialize the server and run it

    if err != nil {
        log.Fatal(err)
    }
}
