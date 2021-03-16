package main

import (
	"log"
	"rabbitmq-go-example/amqp"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("invalid .env")
	}

	// create amqp instance
	amqpInstance := amqp.NewConnection()

	// declare exchange
	amqpInstance.DeclareExchange("EXCHANGE")

	// declare queue
	amqpInstance.DeclareQueue("QUEUE")

	// binding exchange to queue
	amqpInstance.BindingExchangeToQueue("EXCHANGE", "QUEUE", "")

	// create Consumer
	amqpInstance.ReceiveMessage("QUEUE", "")
}
