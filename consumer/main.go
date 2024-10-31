package main

import (
	"encoding/json"
	"fmt"
	"gonext/core"
	"gonext/domain"
	"log"
	"os"
	"time"
)

func main() {

	app := core.NewApplication()

	postgresConn, err := app.Postgres.Connect()
	if err != nil {
		log.Fatalf("Could not connect to postgres: %v", err)
	}

	whatsMeowConn, err := app.WhatsMeowDB.Connect()
	if err != nil {
		log.Fatalf("Could not connect to whatsmeow: %v", err)
	}

	redisConn, err := app.Redis.Connect()
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
	}
	defer redisConn.Close()

	rabbitMqConn, err := app.RabbitMQ.Connect()
	if err != nil {
		log.Fatalf("Could not connect to rabbitmq: %v", err)
	}
	defer rabbitMqConn.Close()

	msgs, err := rabbitMqConn.Consume(os.Getenv("RABBITMQ_QUEUE"))
	if err != nil {
		log.Fatalf("Could not consume messages: %v", err)
	}

	sendMessage := domain.SendMessage{
		WhatsAppRepository: domain.WhatsAppRepository{
			WhatsMeowDB: whatsMeowConn,
			DB:          postgresConn,
		},
	}

	go func() {
		for msg := range msgs {

			var incomingMsg domain.Message
			err := json.Unmarshal(msg.Body, &incomingMsg)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			cache, err := app.Redis.Get(incomingMsg.JID)
			if err != nil && err.Error() != "chave não encontrada" {
				log.Printf("Error getting cache: %v", err)
				msg.Nack(false, true)
				continue
			}

			fmt.Println("cache", cache)

			if cache != "" {
				log.Printf("Message already sent: %v", incomingMsg)
				msg.Nack(false, true)
				continue
			}

			app.Redis.Set(incomingMsg.JID, "true", 0)

			err = sendMessage.Send(&incomingMsg)
			if err != nil {
				app.Redis.Del(incomingMsg.JID)
				log.Printf("Error sending message: %v", err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)

			app.Redis.Del(incomingMsg.JID)

			fmt.Println("incomingMsg", incomingMsg)

		}
	}()

	time.Sleep(5 * time.Second)
	select {}
}
