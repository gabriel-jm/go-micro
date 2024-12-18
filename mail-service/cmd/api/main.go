package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Mailer Mail
}

const webPort = "8000"

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail serving on port:", app)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
