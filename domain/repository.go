package domain

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/segmentio/ksuid"

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

/*
Método CreateAccount cria uma nova conta no banco de dados Postgres.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
- name: Nome do workspace.
- origin: Origem da conta.
- externalId: Identificador externo da conta.
Retorna:
- O ID da conta criada e um erro, se houver.
*/
func (r WhatsAppRepository) CreateAccount(ctx context.Context, name string, origin string, externalId string) (string, error) {

	id := ksuid.New().String()

	query := `
        INSERT INTO accounts (id, name, origin, "externalId", "updatedAt", "createdAt")
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id
    `
	row := r.DB.QueryRowContext(ctx, query, id, name, origin, externalId)
	if row.Err() != nil {
		return "", row.Err()
	}

	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r WhatsAppRepository) CreateSession(ctx context.Context, sessionId string, accountID string, qr string) (string, string, error) {

	var response struct {
		ID       string `json:"id"`
		AuthCode string `json:"authCode"`
	}

	query := `
	INSERT INTO sessions (id, "accountId", "authCode", status, "createdAt", "updatedAt")
	VALUES ($1, $2, $3, $4, NOW(), NOW())
	RETURNING id, "authCode"
	`

	row := r.DB.QueryRowContext(ctx, query, sessionId, accountID, qr, "PENDING")
	if row.Err() != nil {
		return "", "", row.Err()
	}

	err := row.Scan(&response.ID, &response.AuthCode)
	if err != nil {
		return "", "", err
	}

	return response.ID, response.AuthCode, nil

}

func (r WhatsAppRepository) UpdateSession(ctx context.Context, sessionID string, status string, jid string) (err error) {

	jsonJid := map[string]interface{}{
		"phoneId": jid,
	}

	sessionInfoJSON, err := json.Marshal(jsonJid)
	if err != nil {
		return err
	}

	query := `
		UPDATE sessions
		SET status = $1, "sessionInfo" = $2, "readyAt" = $3, "updatedAt" = NOW()
		WHERE id = $4
	`
	row := r.DB.QueryRowContext(ctx, query, status, sessionInfoJSON, "now()", sessionID)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (r WhatsAppRepository) GetSessionByID(ctx context.Context, id string) (res GetSessionByIDResponse, err error) {
	query := `
        SELECT 
            id,
            "accountId",
            "sessionInfo",
            status,
            "authCode",
            "readyAt",
            "failureReason",
            "createdAt",
            "updatedAt"
        FROM sessions
        WHERE id = $1
    `
	row := r.DB.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return GetSessionByIDResponse{}, row.Err()
	}

	var sessionInfoJSON []byte
	res = GetSessionByIDResponse{}
	err = row.Scan(&res.ID, &res.AccountID, &sessionInfoJSON, &res.Status, &res.AuthCode, &res.ReadyAt, &res.FailureReason, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return GetSessionByIDResponse{}, err
	}

	err = json.Unmarshal(sessionInfoJSON, &res.SessionInfo)
	if err != nil {
		return GetSessionByIDResponse{}, err
	}

	return res, nil
}

func (r WhatsAppRepository) DeleteSession(ctx context.Context, id string) error {
	query := `
		DELETE FROM sessions
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r WhatsAppRepository) DeleteAccount(ctx context.Context, id string) error {
	query := `
        DELETE FROM accounts
        WHERE id = $1
    `

	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
