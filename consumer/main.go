package main

import (
	"encoding/json"
	"gozap/core"
	"gozap/domain"
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

	err = app.Messenger.Connect()
	if err != nil {
		log.Fatalf("Could not connect to messenger module: %v", err)
	}
	defer app.Messenger.Close()

	sendMessage := domain.SendMessage{
		WhatsAppRepository: domain.WhatsAppRepository{
			WhatsMeowDB: whatsMeowConn,
			DB:          postgresConn,
		},
	}

	sessionManager := core.NewSessionManager()

	handler := func(msg core.Message) {

		ProcessMessage(sendMessage, sessionManager, msg)
	}

	err = app.Messenger.Consume(os.Getenv("QUEUE_MESSAGE"), handler)
	if err != nil {
		log.Fatalf("Could not consume messages: %v", err)
	}

	time.Sleep(5 * time.Second)
	select {}
}

func ProcessMessage(sendMessage domain.SendMessage, sessionManager *core.SessionManager, msg core.Message) {

	var incomingMsg domain.Message
	err := json.Unmarshal(msg.Body, &incomingMsg)
	if err != nil {
		msg.Nack()
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	if !sessionManager.StartSession(incomingMsg.SessionId) {
		log.Printf("Session %s already in process. Skipping message: %s", incomingMsg.SessionId, incomingMsg.Message)
		msg.Nack()
		return
	}

	err = sendMessage.Send(&domain.Message{
		SessionId: incomingMsg.SessionId,
		To:        incomingMsg.To,
		Message:   incomingMsg.Message,
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
		msg.Nack()
		return
	}

	log.Printf("Processing message for session %s: %s", incomingMsg.SessionId, incomingMsg.Message)
	msg.Ack()
	sessionManager.EndSession(incomingMsg.SessionId)
}
