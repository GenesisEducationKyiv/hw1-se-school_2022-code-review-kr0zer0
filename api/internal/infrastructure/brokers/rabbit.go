package brokers

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BrokerWriter struct {
	connection *amqp.Connection
}

func NewBrokerWriter(connection *amqp.Connection) *BrokerWriter {
	return &BrokerWriter{connection: connection}
}

func (w BrokerWriter) Write(bytes []byte) (int, error) {

	ch, err := w.connection.Channel()
	if err != nil {
		return -1, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"logs", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return -1, err
	}

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

	return len(bytes), nil
}
