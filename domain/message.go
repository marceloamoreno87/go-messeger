package domain

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

type Message struct {
	JID          string `json:"jid"`
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
	if err != nil || deviceStore == nil {
		return err
	}

	client := whatsmeow.NewClient(deviceStore, nil)
	if err = client.Connect(); err != nil {
		return err
	}
	_, err = client.SendMessage(
		context.Background(),
		JID,
		&waProto.Message{
			Conversation: proto.String(message.GetMessage()),
		})

	if err != nil {
		return err
	}

	return nil
}
