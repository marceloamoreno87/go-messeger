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
	JID     string `json:"jid"`
	To      string `json:"to"`
	Message string `json:"message"`
}

func (m *Message) GetMessage() string {
	return m.Message
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
