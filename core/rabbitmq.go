package core

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Config struct {
	URL string
}

type RabbitMQClient struct {
	Config *Config
	Conn   *amqp.Connection
	Ch     *amqp.Channel
}

func (client *RabbitMQClient) Connect() (rabbit *RabbitMQClient, err error) {

	conn, err := amqp.Dial(client.Config.URL)
	if err != nil {
		return client, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	client.Conn = conn

	ch, err := client.Conn.Channel()
	if err != nil {
		return client, fmt.Errorf("failed to open a channel: %w", err)
	}
	client.Ch = ch

	return client, nil

}

func (client *RabbitMQClient) Publish(queueName string, body string) (rabbit *RabbitMQClient, err error) {
	q, err := client.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return client, fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = client.Ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return client, fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Message sent: %s", body)
	return client, nil
}

func (client *RabbitMQClient) Consume(queueName string) (<-chan amqp.Delivery, error) {
	q, err := client.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	msgs, err := client.Ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return msgs, nil
}

func (client *RabbitMQClient) Close() {
	if client.Ch != nil {
		client.Ch.Close()
	}
	if client.Conn != nil {
		client.Conn.Close()
	}
}
