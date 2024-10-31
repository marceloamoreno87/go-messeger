package domain

import (
	"context"
	"database/sql"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

type WhatsAppRepository struct {
	WhatsMeowDB *sqlstore.Container
	DB          *sql.DB
}

func (r WhatsAppRepository) CreateDeviceWM(ctx context.Context) (device *store.Device) {
	return r.WhatsMeowDB.NewDevice()
}

func (r WhatsAppRepository) DeleteDeviceWM(ctx context.Context, device *store.Device) (err error) {
	return r.WhatsMeowDB.DeleteDevice(device)
}

func (r WhatsAppRepository) FindDeviceWM(ctx context.Context, JID types.JID) (device *store.Device, err error) {
	return r.WhatsMeowDB.GetDevice(JID)
}

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

func (r WhatsAppRepository) DeleteDevicePG(ctx context.Context, phone string) error {
	query := `
	DELETE FROM devices
	WHERE phone_number = $1
	`
	row := r.DB.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}
