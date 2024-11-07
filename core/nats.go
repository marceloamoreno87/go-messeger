package core

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
)

/*
Definição de variáveis de erro específicas para operações com o NATS.
Essas variáveis são usadas para fornecer mensagens de erro detalhadas.
*/
var (
	ErrNatsConnectionFailed = errors.New("nats.connection_failed: failed to connect to NATS")
	ErrNatsPublish          = errors.New("nats.publish_failed: failed to publish a message")
	ErrNatsSubscribe        = errors.New("nats.subscribe_failed: failed to subscribe to a subject")
	ErrNatsCloseConnection  = errors.New("nats.close_connection_failed: failed to close connection")
)

type NatsInterface interface {
}

/*
Estrutura NatsMessenger que mantém a configuração e conexão do NATS.
*/
type NatsMessenger struct {
	URL  string
	conn *nats.Conn
	js   nats.JetStreamContext
}

/*
Método Connect estabelece uma conexão com o NATS.
Retorna um ponteiro para o NatsMessenger e um erro, se houver.
*/
func (client *NatsMessenger) Connect() error {
	nc, err := nats.Connect(client.URL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNatsConnectionFailed, err)
	}
	client.conn = nc
	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNatsConnectionFailed, err)
	}
	client.js = js
	return nil
}

/*
Método Publish publica uma mensagem no assunto especificado.
Retorna um erro, se houver.
*/
func (client *NatsMessenger) Publish(subject string, body []byte) error {
	_, err := client.js.Publish(subject, body)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNatsPublish, err)
	}
	return nil
}

/*
Método Subscribe consome mensagens do assunto especificado.
Retorna um canal de mensagens (<-chan *nats.Msg) e um erro, se houver.
*/
func (client *NatsMessenger) Consume(subject string, handler func(Message)) error {
	_, err := client.js.Subscribe(subject, func(msg *nats.Msg) {
		message := Message{
			Body: msg.Data,
			Ack: func() error {
				return msg.Ack()
			},
			Nack: func() error {
				return msg.Nak()
			},
		}
		handler(message)
	})
	return err
}

/*
Método Close fecha a conexão do NATS.
Retorna um erro se houver falha ao fechar a conexão.
*/
func (client *NatsMessenger) Close() error {
	if client.conn != nil {
		client.conn.Close()
	}
	return nil
}
