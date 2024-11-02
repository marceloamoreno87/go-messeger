package domain

import (
	"context"
	"database/sql"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

/*
Estrutura WhatsAppRepository que contém as conexões para o banco de dados WhatsMeow e Postgres.
Esta estrutura é responsável por realizar operações no banco de dados relacionadas ao WhatsApp.
*/
type WhatsAppRepository struct {
	WhatsMeowDB *sqlstore.Container
	DB          *sql.DB
}

/*
Método CreateDeviceWM cria um novo dispositivo no banco de dados WhatsMeow.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
Retorna:
- Um ponteiro para a estrutura store.Device.
*/
func (r WhatsAppRepository) CreateDeviceWM(ctx context.Context) (device *store.Device) {
	return r.WhatsMeowDB.NewDevice()
}

/*
Método DeleteDeviceWM deleta um dispositivo do banco de dados WhatsMeow.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
- device: Ponteiro para a estrutura store.Device a ser deletada.
Retorna:
- Um erro, se houver.
*/
func (r WhatsAppRepository) DeleteDeviceWM(ctx context.Context, device *store.Device) (err error) {
	return r.WhatsMeowDB.DeleteDevice(device)
}

/*
Método FindDeviceWM encontra um dispositivo no banco de dados WhatsMeow pelo JID.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
- JID: Identificador do dispositivo a ser encontrado.
Retorna:
- Um ponteiro para a estrutura store.Device e um erro, se houver.
*/
func (r WhatsAppRepository) FindDeviceWM(ctx context.Context, sessionID string) (device *store.Device, err error) {

	query := `
		SELECT 
		jid
		FROM whatsmeow_device 
		WHERE registration_id = $1
		`
	row := r.DB.QueryRowContext(ctx, query, sessionID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var jid string
	err = row.Scan(&jid)
	if err != nil {
		return nil, err
	}

	device, err = r.WhatsMeowDB.GetDevice(NumberToJID(jid))
	if err != nil {
		return nil, err
	}

	return device, nil
}

/*
Método UpdateOrCreatePG atualiza ou cria um dispositivo no banco de dados Postgres.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
- JID: Identificador do dispositivo.
- phoneNumber: Número de telefone associado ao dispositivo.
Retorna:
- Um erro, se houver.
*/
func (r WhatsAppRepository) UpdateOrCreatePG(ctx context.Context, JID string, phoneNumber string) error {
	query := `
        INSERT INTO devices (jid, phone_number)
        VALUES ($1, $2)
        ON CONFLICT (phone_number) DO UPDATE
        SET jid = EXCLUDED.jid, updated_at = CURRENT_TIMESTAMP
    `

	row := r.DB.QueryRowContext(ctx, query, JID, phoneNumber)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

/*
Método DeleteDevicePG deleta um dispositivo do banco de dados Postgres pelo número de telefone.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
- phone: Número de telefone do dispositivo a ser deletado.
Retorna:
- Um erro, se houver.
*/
func (r WhatsAppRepository) DeleteDevicePG(ctx context.Context, phone string) error {
	query := `
        DELETE FROM devices
        WHERE phone_number = $1
    `
	row := r.DB.QueryRowContext(ctx, query, phone)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}
