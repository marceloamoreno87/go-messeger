package domain

import (
	"strings"

	"go.mau.fi/whatsmeow/types"
)

func NumberToJID(JID string) types.JID {
	parts := strings.Split(JID, "@")
	return types.NewJID(parts[0], parts[1])
}

func JIDToNumber(jid types.JID) string {
	return jid.User
}
