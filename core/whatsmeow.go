package core

import (
	"database/sql"

	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type WhatsMeowDB struct {
	DSN string
}

func (wm *WhatsMeowDB) Connect() (postgresDb *sqlstore.Container, err error) {

	db, err := sql.Open("postgres", wm.DSN)
	if err != nil {
		log.Error().Msg("Error trying to connect to database: " + err.Error())
		return
	}
	defer db.Close()

	postgresDb, err = sqlstore.New("postgres", wm.DSN, nil)
	if err != nil {
		log.Error().Msg("Err to create container with DB" + err.Error())
		return
	}

	return
}
