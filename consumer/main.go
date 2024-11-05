package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gonext/core"
	"gonext/domain"
	"log"
	"os"
	"time"
)

func main() {

	/*
	   Cria uma nova instância da aplicação.
	   A função NewApplication inicializa todas as dependências necessárias.
	*/
	app := core.NewApplication()

	/*
	   Conecta ao banco de dados Postgres.
	   Se a conexão falhar, o programa será encerrado com uma mensagem de erro.
	*/
	postgresConn, err := app.Postgres.Connect()
	if err != nil {
		log.Fatalf("Could not connect to postgres: %v", err)
	}

	/*
	   Conecta ao banco de dados WhatsMeow.
	   Se a conexão falhar, o programa será encerrado com uma mensagem de erro.
	*/
	whatsMeowConn, err := app.WhatsMeowDB.Connect()
	if err != nil {
		log.Fatalf("Could not connect to whatsmeow: %v", err)
	}

	/*
	   Conecta ao Redis.
	   Se a conexão falhar, o programa será encerrado com uma mensagem de erro.
	*/
	redisConn, err := app.Redis.Connect()
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
	}
	defer redisConn.Close()

	/*
	   Conecta ao RabbitMQ.
	   Se a conexão falhar, o programa será encerrado com uma mensagem de erro.
	*/
	rabbitMqConn, err := app.RabbitMQ.Connect()
	if err != nil {
		log.Fatalf("Could not connect to rabbitmq: %v", err)
	}
	defer rabbitMqConn.Close()

	/*
	   Consome mensagens da fila RabbitMQ.
	   Se a operação falhar, o programa será encerrado com uma mensagem de erro.
	*/
	msgs, err := rabbitMqConn.Consume(os.Getenv("RABBITMQ_QUEUE"))
	if err != nil {
		log.Fatalf("Could not consume messages: %v", err)
	}

	/*
	   Cria uma instância de SendMessage com o repositório WhatsApp.
	   O repositório é responsável por interagir com o banco de dados WhatsMeow e Postgres.
	*/
	sendMessage := domain.SendMessage{
		WhatsAppRepository: domain.WhatsAppRepository{
			WhatsMeowDB: whatsMeowConn,
			DB:          postgresConn,
		},
	}

	/*
	   Inicia uma goroutine para processar mensagens da fila RabbitMQ.
	   A goroutine continua processando mensagens até que o canal seja fechado.
	*/
	go func() {
		for msg := range msgs {

			fmt.Println("START", string(msg.Body))

			/*
			   Deserializa a mensagem recebida do formato JSON para a estrutura Message.
			   Se a deserialização falhar, a mensagem é ignorada e o erro é registrado.
			*/
			var incomingMsg domain.Message
			err := json.Unmarshal(msg.Body, &incomingMsg)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			/*
			   Verifica se a mensagem já foi processada anteriormente consultando o cache Redis.
			   Se houver um erro ao consultar o cache, a mensagem é reenviada para a fila.
			*/
			cache, err := app.Redis.Get(incomingMsg.SessionId)
			if err != nil && !errors.Is(err, core.ErrRedisKeyNotFound) {
				log.Printf("Error getting cache: %v", err)
				msg.Nack(false, true)
				continue
			}

			/*
			   Se o JID está em uso por outro dispositivo, a mensagem é volta pra fila.
			*/
			if cache == "true" {
				log.Printf("Message already sent: %v", incomingMsg)
				msg.Nack(false, true)
				continue
			}

			/*
			   Armazena a mensagem no cache Redis para evitar processamento duplicado.
			*/
			app.Redis.Set(incomingMsg.SessionId, "true", 0)

			_, err = app.Redis.Get(incomingMsg.SessionId)
			if err != nil && !errors.Is(err, core.ErrRedisKeyNotFound) {
				log.Printf("Error getting cache: %v", err)
				msg.Nack(false, true)
				continue
			}

			/*
			   Envia a mensagem usando o serviço SendMessage.
			   Se o envio falhar, a mensagem é removida do cache e reenviada para a fila.
			*/
			err = sendMessage.Send(&incomingMsg)
			if err != nil {
				if err.Error() == "device_not_found" {
					log.Printf("Device not found: %v", incomingMsg)
					app.Redis.Del(incomingMsg.SessionId)
					msg.Ack(false)
					continue
				}
				app.Redis.Del(incomingMsg.SessionId)
				log.Printf("Error sending message: %v", err)
				msg.Nack(false, true)
				continue
			}

			/*
			   Reconhece a mensagem como processada com sucesso.
			*/
			msg.Ack(false)

			/*
			   Remove a mensagem do cache após o processamento bem-sucedido.
			*/
			app.Redis.Del(incomingMsg.SessionId)

			fmt.Println("FINISH", incomingMsg)

		}
	}()

	/*
	   Mantém o programa em execução por 5 segundos para permitir o processamento de mensagens.
	*/
	time.Sleep(5 * time.Second)
	select {}
}
