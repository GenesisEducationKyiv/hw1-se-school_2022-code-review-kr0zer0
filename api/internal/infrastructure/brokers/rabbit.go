package brokers

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type BrokerWriter struct {
}

func (w BrokerWriter) Write(bytes []byte) (int, error) {
	conn, err := amqp.Dial("amqp://test:test@localhost:5672/")
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return -1, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return -1, err
	}

	body := "Hello World!"
	err = ch.PublishWithContext(context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})
	if err != nil {
		return -1, err
	}
	log.Printf(" [x] Sent %s\n", body)
	return len(bytes), nil
}
