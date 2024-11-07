package main

import (
	"fmt"
	"listener/events"
	"log"
	"math"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connect()

	if err != nil {
		log.Fatalln(err)
	}

	defer rabbitConn.Close()

	consumer, err := events.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Fatal(err)
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
