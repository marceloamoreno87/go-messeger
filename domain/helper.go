package domain

import (
	"strings"

	"go.mau.fi/whatsmeow/types"
)

/*
Função NumberToJID converte um número de telefone no formato JID (Jabber ID).
Parâmetros:
- JID: String contendo o número de telefone no formato JID (ex: "12345@s.whatsapp.net").
Retorna:
- Um objeto types.JID contendo o número de telefone e o servidor.
*/
func NumberToJID(JID string) types.JID {
	parts := strings.Split(JID, "@")
	return types.NewJID(parts[0], parts[1])
}

/*
Função JIDToNumber converte um objeto types.JID em uma string contendo apenas o número de telefone.
Parâmetros:
- jid: Objeto types.JID contendo o número de telefone e o servidor.
Retorna:
- Uma string contendo apenas o número de telefone.
*/
func JIDToNumber(jid types.JID) string {
	return jid.User
}
