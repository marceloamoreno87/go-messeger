package core

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

/*
Definição de variáveis de erro específicas para operações com o RabbitMQ.
Essas variáveis são usadas para fornecer mensagens de erro detalhadas.
*/
var (
	ErrRabbitMQConnectionFailed = errors.New("rabbitmq.connection_failed: failed to connect to RabbitMQ")
	ErrRabbitMQChannelFailed    = errors.New("rabbitmq.channel_failed: failed to open a channel")
	ErrRabbitMQQueueDeclare     = errors.New("rabbitmq.queue_declare_failed: failed to declare a queue")
	ErrRabbitMQPublish          = errors.New("rabbitmq.publish_failed: failed to publish a message")
	ErrRabbitMQConsume          = errors.New("rabbitmq.consume_failed: failed to consume messages")
	ErrRabbitMQCloseChannel     = errors.New("rabbitmq.close_channel_failed: failed to close channel")
	ErrRabbitMQCloseConnection  = errors.New("rabbitmq.close_connection_failed: failed to close connection")
)

type RabbitMQMessagerInterface interface {
}

/*
Estrutura RabbitMQMessenger que mantém a configuração, conexão e canal do RabbitMQ.
*/
type RabbitMQMessenger struct {
	URL  string
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

/*
Método Connect estabelece uma conexão com o RabbitMQ e abre um canal.
Retorna um ponteiro para o RabbitMQMessenger e um erro, se houver.
*/
func (client *RabbitMQMessenger) Connect() error {
	conn, err := amqp.Dial(client.URL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRabbitMQConnectionFailed, err)
	}
	client.Conn = conn

	ch, err := client.Conn.Channel()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRabbitMQChannelFailed, err)
	}
	client.Ch = ch

	return nil
}

/*
Método Publish publica uma mensagem na fila especificada.
Retorna um ponteiro para o RabbitMQMessenger e um erro, se houver.
*/
func (client *RabbitMQMessenger) Publish(queueName string, body []byte) error {
	q, err := client.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRabbitMQQueueDeclare, err)
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
		return fmt.Errorf("%w: %v", ErrRabbitMQPublish, err)
	}

	log.Printf("Message sent: %s", body)
	return nil
}

/*
Método Consume consome mensagens da fila especificada.
Retorna um canal de mensagens (<-chan amqp.Delivery) e um erro, se houver.
*/
func (client *RabbitMQMessenger) Consume(queueName string, handler func(Message)) error {
	q, err := client.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRabbitMQQueueDeclare, err)
	}

	msgs, err := client.Ch.Consume(
		q.Name,
		"",
		false, // autoAck set to false for manual acknowledgment
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRabbitMQConsume, err)
	}

	go func() {
		for d := range msgs {
			msg := Message{
				Body: d.Body,
				Ack: func() error {
					return d.Ack(false)
				},
				Nack: func() error {
					return d.Nack(false, true)
				},
			}
			handler(msg)
		}
	}()

	return nil
}

/*
Método Close fecha o canal e a conexão do RabbitMQ.
Retorna um erro se houver falha ao fechar o canal ou a conexão.
*/
func (client *RabbitMQMessenger) Close() error {
	var errStrings []string

	if err := client.Ch.Close(); err != nil {
		errStrings = append(errStrings, fmt.Sprintf("%v: %v", ErrRabbitMQCloseChannel, err))
	}
	if err := client.Conn.Close(); err != nil {
		errStrings = append(errStrings, fmt.Sprintf("%v: %v", ErrRabbitMQCloseConnection, err))
	}

	if len(errStrings) > 0 {
		return errors.New(fmt.Sprintf("multiple errors occurred: %s", errStrings))
	}
	return nil
}
