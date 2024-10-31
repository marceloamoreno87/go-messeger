package main

import (
	"fmt"
	"gonext/core"
	"gonext/domain"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	
	app := core.NewApplication()

	postgresConn, err := app.Postgres.Connect()
	if err != nil {
		log.Fatalf("Could not connect to postgres: %v", err)
	}
	app.Postgres.RunMigrate()

	whatsMeowConn, err := app.WhatsMeowDB.Connect()
	if err != nil {
		log.Fatalf("Could not connect to whatsmeow: %v", err)
	}

	rabbitMqConn, err := app.RabbitMQ.Connect()
	if err != nil {
		log.Fatalf("Could not connect to rabbitmq: %v", err)
	}
	defer rabbitMqConn.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler := domain.WhatsAppHandler{
		WhatsAppService: domain.WhatsAppService{
			RabbitMQService: rabbitMqConn,
			WhatsAppRepository: domain.WhatsAppRepository{
				WhatsMeowDB: whatsMeowConn,
				DB:          postgresConn,
			},
		},
	}

	r.Get("/connect", handler.Connect)
	r.Post("/validate", handler.Validate)
	r.Post("/send", handler.Send)

	fmt.Println("API running on port " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
