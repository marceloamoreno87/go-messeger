package domain

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type Message struct {
	JID          string `json:"jid"`
	To           string `json:"to"`
	Title        string `json:"title"`
	Service      string `json:"service"`
	When         string `json:"when"`
	Duration     string `json:"duration"`
	Professional string `json:"professional"`
	Code         string `json:"code"`
	Footer       string `json:"footer"`
}

func (m *Message) GetMessage() string {
	return m.Title + "\n" +
		"Serviço: " + m.Service + "\n" +
		"Quando: " + m.When + "\n" +
		"Duração: " + m.Duration + "\n" +
		"Profissional: " + m.Professional + "\n" +
		"Código: " + m.Code + "\n" +
		"Rodapé: " + m.Footer
}

type SendMessage struct {
	WhatsAppRepository WhatsAppRepository
}

func (s *SendMessage) Send(message *Message) (err error) {
	JID := NumberToJID(message.JID)
	deviceStore, err := s.WhatsAppRepository.FindDeviceWM(context.Background(), JID)
	if err != nil {
		return err
	}
	fmt.Println(deviceStore)

	client := whatsmeow.NewClient(deviceStore, nil)
	if err = client.Connect(); err != nil {
		return err
	}

	TO := types.JID{
		Server: "s.whatsapp.net",
		User:   message.To,
	}

	_, err = client.SendMessage(
		context.Background(),
		TO,
		&waProto.Message{
			Conversation: proto.String(message.GetMessage()),
		})

	if err != nil {
		return err
	}

	return nil
}
