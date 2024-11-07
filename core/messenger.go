package core

import (
	"os"
)

type Message struct {
	Body []byte
	Ack  func() error
	Nack func() error
}

type MessengerInterface interface {
	Connect() error
	Publish(queueName string, body []byte) error
	Consume(queueName string, handler func(Message)) error
	Close() error
}

type DriverMessage int

const (
	RabbitMQ DriverMessage = iota
	Nats
)

func (d DriverMessage) String() string {
	return [...]string{"RabbitMQ", "Nats"}[d]
}

func ParseDriverMessage(s string) DriverMessage {
	switch s {
	case "RabbitMQ":
		return RabbitMQ
	case "Nats":
		return Nats
	default:
		return RabbitMQ
	}
}

type Messenger struct {
	DriverMessage DriverMessage
}

func NewMessenger(driverMessage DriverMessage) MessengerInterface {
	switch driverMessage {
	case RabbitMQ:
		return &RabbitMQMessenger{
			URL: os.Getenv("RABBITMQ_DSN"),
		}
	case Nats:
		return &NatsMessenger{
			URL: os.Getenv("NATS_DSN"),
		}
	default:
		return &RabbitMQMessenger{
			URL: os.Getenv("RABBITMQ_DSN"),
		}
	}
}
