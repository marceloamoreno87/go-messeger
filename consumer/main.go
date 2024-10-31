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

	whatsMeowDB := core.WhatsMeowDB{
		DSN: os.Getenv("POSTGRES_DSN"),
	}
	connWhatsMeow, err := whatsMeowDB.DbSqlConnect()
	if err != nil {
		log.Fatal(err)
	}

	database := core.Postgres{
		DSN: os.Getenv("POSTGRES_DSN"),
	}

	conn, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	database.RunMigrate()

	client := &core.RabbitMQClient{
		Config: &core.Config{
			URL: os.Getenv("RABBITMQ_DSN"),
		},
	}

	sendMessage := domain.SendMessage{
		WhatsAppRepository: domain.WhatsAppRepository{
			WhatsMeowDB: connWhatsMeow,
			DB:          conn,
		},
	}

	redisClient := core.RedisClient{
		DSN: os.Getenv("REDIS_DSN"),
	}

	defer redisClient.Close()

	client, err = client.Connect()
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer client.Close()

	msgs, err := client.Consume(os.Getenv("RABBITMQ_QUEUE"))
	if err != nil {
		log.Fatalf("Could not consume messages: %v", err)
	}

	go func() {
		for msg := range msgs {

			var incomingMsg domain.Message
			err := json.Unmarshal(msg.Body, &incomingMsg)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			cache, err := redisClient.Get(incomingMsg.JID)
			if err != nil {
				log.Printf("Error getting cache: %v", err)
				msg.Nack(false, true)
				continue
			}

			if cache != "" {
				log.Printf("Message already sent: %v", incomingMsg)
				msg.Nack(false, true)
				continue
			}

			redisClient.Set(incomingMsg.JID, "true", 0)

			err = sendMessage.Send(&incomingMsg)
			if err != nil {
				redisClient.Del(incomingMsg.JID)
				log.Printf("Error sending message: %v", err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)

			redisClient.Del(incomingMsg.JID)

			fmt.Println("incomingMsg", incomingMsg)

		}
	}()

	time.Sleep(5 * time.Second)
	select {}
}
