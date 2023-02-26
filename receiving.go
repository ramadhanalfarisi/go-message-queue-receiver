package main

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type send struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (s *send) connectRabbitMq() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	s.Connection = conn
}

func (s *send) channelRabbitMq() {
	ch, err := s.Connection.Channel()
	failOnError(err, "Failed to open a channel")
	s.Channel = ch
}

func main() {
	var s send
	s.connectRabbitMq()
	defer s.Connection.Close()
	s.channelRabbitMq()
	defer s.Channel.Close()
	q, err := s.Channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := s.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
