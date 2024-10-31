package core

import "os"

type Application struct {
	WhatsMeowDB WhatsMeowDB
	Postgres    Postgres
	RabbitMQ    RabbitMQClient
	Redis       RedisClient
}

func NewApplication() *Application {
	return &Application{
		WhatsMeowDB: WhatsMeowDB{
			DSN: os.Getenv("POSTGRES_DSN"),
		},
		Postgres: Postgres{
			DSN: os.Getenv("POSTGRES_DSN"),
		},
		RabbitMQ: RabbitMQClient{
			Config: &Config{
				URL: os.Getenv("RABBITMQ_DSN"),
			},
		},
		Redis: RedisClient{
			DSN: os.Getenv("REDIS_DSN"),
		},
	}
}
