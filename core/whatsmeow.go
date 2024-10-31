package core

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

/*
Definição de variáveis de erro específicas para operações com o banco de dados WhatsMeow.
Essas variáveis são usadas para fornecer mensagens de erro detalhadas.
*/
var (
	ErrDBConnectionFailedWM = errors.New("whatsmeow.db_connection_failed: failed to connect to the database")
	ErrContainerCreation    = errors.New("whatsmeow.container_creation_failed: failed to create container with DB")
)

/*
Estrutura WhatsMeowDB que contém a string de conexão DSN.
DSN (Data Source Name) é usado para conectar ao banco de dados WhatsMeow.
*/
type WhatsMeowDB struct {
	DSN string
}

/*
Método Connect estabelece uma conexão com o banco de dados WhatsMeow.
Retorna um ponteiro para a estrutura sqlstore.Container e um erro, se houver.
*/
func (wm *WhatsMeowDB) Connect() (postgresDb *sqlstore.Container, err error) {
	/*
	   Abre uma conexão com o banco de dados Postgres usando a string de conexão DSN.
	   Se a conexão falhar, um erro detalhado é registrado e retornado.
	*/
	db, err := sql.Open("postgres", wm.DSN)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrDBConnectionFailedWM, err)
		log.Error().Msg(err.Error())
		return
	}
	defer db.Close()

	/*
	   Cria um novo container sqlstore para o banco de dados WhatsMeow.
	   Se a criação do container falhar, um erro detalhado é registrado e retornado.
	*/
	postgresDb, err = sqlstore.New("postgres", wm.DSN, nil)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrContainerCreation, err)
		log.Error().Msg(err.Error())
		return
	}

	return
}
