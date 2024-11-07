package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const webPort = "8000"

type Config struct {
	Rabbit *amqp091.Connection
}

func main() {
	rabbitConn, err := connect()

	if err != nil {
		log.Fatalln(err)
	}

	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp091.Connection, error) {
	var counts int8
	var backOffTime = 1 * time.Second
	var connection *amqp091.Connection

	for {
		c, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672")

		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		log.Println("Backing off...")
		time.Sleep(backOffTime)

		backOffTime = time.Duration(math.Pow(float64(counts), 2)) * time.Second

		continue
	}

	return connection, nil
}
