package amqp

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ interface {
	DeclareExchange(exchangeName string)
	DeclareQueue(queueName string)
	BindingExchangeToQueue(exchangeName string, queueName string, route string)
	SendMessage(exchangeName string, route string, message string)
	ReceiveMessage(queueName string, route string)
}

type rabbitMQ struct {
	connection *amqp.Connection
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewConnection() *rabbitMQ {
	uri := os.Getenv("RABBITMQ_URI")

	if uri == "" {
		log.Fatal("RabbitMQ URI is invalid")
	}

	connection, err := amqp.Dial(uri)
	failOnError(err, "Failed connect to RabbiqMQ Server")

	log.Println("Connected to RabbitMQ")

	return &rabbitMQ{connection}
}

func (rabbitmq *rabbitMQ) DeclareExchange(exchangeName string) {
	ch, err := rabbitmq.connection.Channel()
	failOnError(err, "Failed to open channel")

	err = ch.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	log.Printf("RabbitMQ: %s exchange created", exchangeName)
}

func (rabbitmq *rabbitMQ) DeclareQueue(queueName string) {
	ch, err := rabbitmq.connection.Channel()
	failOnError(err, "Failed to open channel")

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)

	failOnError(err, "Failed to declare a queue")
}

func (rabbitmq *rabbitMQ) BindingExchangeToQueue(exchangeName string, queueName string, route string) {
	ch, err := rabbitmq.connection.Channel()
	failOnError(err, "Failed to open channel")

	err = ch.QueueBind(
		queueName,
		route,
		exchangeName,
		false,
		nil,
	)

	failOnError(err, "Failed to binding queue with exchange")

	log.Printf("RabbitMQ: Binding %s to %s successful", exchangeName, queueName)
}

func (rabbitmq *rabbitMQ) SendMessage(exchangeName string, route string, message string) {
	ch, err := rabbitmq.connection.Channel()

	failOnError(err, "Failed to open channel")

	content := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	err = ch.Publish(
		exchangeName,
		route,
		false,
		false,
		content,
	)

	failOnError(err, "Failed to publish content")

	log.Printf("RabbitMQ: published message to %s", exchangeName)
}

func (rabbitmq *rabbitMQ) ReceiveMessage(queueName string, route string) {
	ch, err := rabbitmq.connection.Channel()
	failOnError(err, "Failed to open Channel")
	msgs, err := ch.Consume(
		queueName,
		route,
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register consumer")
	log.Printf("RabbitMQ: Consumer created - %s -> %s", queueName, route)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("[x] New Message %s - %s : %s", queueName, route, d.Body)
		}
	}()
	<-forever
}
