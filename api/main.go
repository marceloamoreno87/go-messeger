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

	err = database.RunMigrate()
	if err != nil {
		log.Fatal(err)
	}

	rabbitConn := core.RabbitMQClient{
		Config: &core.Config{
			URL: os.Getenv("RABBITMQ_DSN"),
		},
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler := domain.WhatsAppHandler{
		WhatsAppService: domain.WhatsAppService{
			RabbitMQService: rabbitConn,
			WhatsAppRepository: domain.WhatsAppRepository{
				WhatsMeowDB: connWhatsMeow,
				DB:          conn,
			},
		},
	}

	r.Get("/connect", handler.Connect)
	r.Post("/validate", handler.Validate)
	r.Post("/send", handler.Send)

	fmt.Println("API running on port " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
