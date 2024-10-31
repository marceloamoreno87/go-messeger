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

/*
Estrutura Config que contém a URL de conexão para o RabbitMQ.
*/
type Config struct {
	URL string
}

/*
Estrutura RabbitMQClient que mantém a configuração, conexão e canal do RabbitMQ.
*/
type RabbitMQClient struct {
	Config *Config
	Conn   *amqp.Connection
	Ch     *amqp.Channel
}

/*
Método Connect estabelece uma conexão com o RabbitMQ e abre um canal.
Retorna um ponteiro para o RabbitMQClient e um erro, se houver.
*/
func (client *RabbitMQClient) Connect() (*RabbitMQClient, error) {
	conn, err := amqp.Dial(client.Config.URL)
	if err != nil {
		return client, fmt.Errorf("%w: %v", ErrRabbitMQConnectionFailed, err)
	}
	client.Conn = conn

	ch, err := client.Conn.Channel()
	if err != nil {
		return client, fmt.Errorf("%w: %v", ErrRabbitMQChannelFailed, err)
	}
	client.Ch = ch

	return client, nil
}

/*
Método Publish publica uma mensagem na fila especificada.
Retorna um ponteiro para o RabbitMQClient e um erro, se houver.
*/
func (client *RabbitMQClient) Publish(queueName string, body string) (*RabbitMQClient, error) {
	q, err := client.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return client, fmt.Errorf("%w: %v", ErrRabbitMQQueueDeclare, err)
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
		return client, fmt.Errorf("%w: %v", ErrRabbitMQPublish, err)
	}

	log.Printf("Message sent: %s", body)
	return client, nil
}

/*
Método Consume consome mensagens da fila especificada.
Retorna um canal de mensagens (<-chan amqp.Delivery) e um erro, se houver.
*/
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
		return nil, fmt.Errorf("%w: %v", ErrRabbitMQQueueDeclare, err)
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
		return nil, fmt.Errorf("%w: %v", ErrRabbitMQConsume, err)
	}

	return msgs, nil
}

/*
Método Close fecha o canal e a conexão do RabbitMQ.
Retorna um erro se houver falha ao fechar o canal ou a conexão.
*/
func (client *RabbitMQClient) Close() error {
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
