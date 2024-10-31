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
	   Executa as migrações do banco de dados.
	   Isso garante que o esquema do banco de dados esteja atualizado.
	*/
	app.Postgres.RunMigrate()

	/*
	   Conecta ao banco de dados WhatsMeow.
	   Se a conexão falhar, o programa será encerrado com uma mensagem de erro.
	*/
	whatsMeowConn, err := app.WhatsMeowDB.Connect()
	if err != nil {
		log.Fatalf("Could not connect to whatsmeow: %v", err)
	}

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
	   Cria um novo roteador usando o pacote chi.
	   O middleware Logger é usado para registrar todas as solicitações HTTP.
	*/
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	/*
	   Cria um novo manipulador para o serviço WhatsApp.
	   O manipulador é responsável por lidar com as solicitações HTTP relacionadas ao WhatsApp.
	*/
	handler := domain.WhatsAppHandler{
		WhatsAppService: domain.WhatsAppService{
			RabbitMQService: rabbitMqConn,
			WhatsAppRepository: domain.WhatsAppRepository{
				WhatsMeowDB: whatsMeowConn,
				DB:          postgresConn,
			},
		},
	}

	/*
	   Define as rotas HTTP e os manipuladores correspondentes.
	   /connect: Manipulador para conectar ao serviço WhatsApp.
	   /validate: Manipulador para validar dados.
	   /send: Manipulador para enviar mensagens.
	*/
	r.Get("/connect", handler.Connect)
	r.Post("/validate", handler.Validate)
	r.Post("/send", handler.Send)

	/*
	   Inicia o servidor HTTP na porta especificada.
	   A porta é obtida a partir da variável de ambiente PORT.
	*/
	fmt.Println("API running on port " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
