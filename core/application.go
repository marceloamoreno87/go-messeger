package core

import (
	"os"
)

/*
Estrutura Application que contém todas as dependências necessárias para a aplicação.
Inclui conexões para WhatsMeowDB, Postgres, RabbitMQ e Redis.
*/
type Application struct {
	WhatsMeowDB WhatsMeowDB
	Postgres    Postgres
	Messenger   MessengerInterface
	Redis       RedisClient
}

/*
Função NewApplication cria uma nova instância da estrutura Application.
Inicializa todas as dependências necessárias usando variáveis de ambiente.
Retorna um ponteiro para a estrutura Application.
*/
func NewApplication() *Application {

	driverMessenger := ParseDriverMessage(os.Getenv("MESSENGER_DRIVER"))
	messenger := NewMessenger(driverMessenger)

	return &Application{
		/*
		   Inicializa a conexão com o banco de dados WhatsMeowDB.
		   A string de conexão é obtida da variável de ambiente POSTGRES_DSN.
		*/
		WhatsMeowDB: WhatsMeowDB{
			DSN: os.Getenv("POSTGRES_DSN"),
		},
		/*
		   Inicializa a conexão com o banco de dados Postgres.
		   A string de conexão é obtida da variável de ambiente POSTGRES_DSN.
		*/
		Postgres: Postgres{
			DSN: os.Getenv("POSTGRES_DSN"),
		},
		/*
		   Inicializa a conexão com o RabbitMQ.
		   A URL de conexão é obtida da variável de ambiente RABBITMQ_DSN.
		*/
		Messenger: messenger,
		/*
		   Inicializa a conexão com o Redis.
		   A string de conexão é obtida da variável de ambiente REDIS_DSN.
		*/
		Redis: RedisClient{
			DSN: os.Getenv("REDIS_DSN"),
		},
	}
}
