package domain

import (
	"context"
	"gonext/core"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types/events"
)

type WhatsAppService struct {
	WhatsAppRepository WhatsAppRepository
	RabbitMQService    core.RabbitMQClient
}

func (s WhatsAppService) Connect(ctx context.Context) ([]byte, error) {

	deviceStore := s.WhatsAppRepository.CreateDeviceWM(context.Background())
	client := whatsmeow.NewClient(deviceStore, nil)

	pairSuccessChan := make(chan *events.PairSuccess, 1)
	client.AddEventHandler(func(evt interface{}) {
		if pairSuccess, ok := evt.(*events.PairSuccess); ok {
			pairSuccessChan <- pairSuccess
		}
	})

	go func() {
		select {
		case <-pairSuccessChan:
			err := s.WhatsAppRepository.SaveDevicePG(context.Background(), deviceStore.ID.ADString(), deviceStore.ID.User)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	store.SetOSInfo("Windows", [3]uint32{1, 2, 3})

	qrChan, _ := client.GetQRChannel(context.Background())
	err := client.Connect()
	if err != nil {
		return nil, err
	}
	qr := []byte{}
	for evt := range qrChan {
		if evt.Event == "code" {
			qr, err = qrcode.Encode(evt.Code, qrcode.Medium, 256)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	return qr, nil
}

func (s WhatsAppService) Validate(ctx context.Context, req ValidateRequest) (res ValidateResponse, err error) {
	JID := NumberToJID(req.JID)
	device, err := s.WhatsAppRepository.FindDeviceWM(ctx, JID)
	if err != nil || device == nil {
		return ValidateResponse{
			Active: false,
		}, err
	}

	return ValidateResponse{
		Active: true,
	}, nil

}

func (s WhatsAppService) Send(ctx context.Context, req SendRequest) (res SendResponse, err error) {
	client, err := s.RabbitMQService.Connect()
	if err != nil {
		return SendResponse{
			Sent: false,
		}, err
	}
	defer client.Close()

	client, err = client.Publish(os.Getenv("RABBITMQ_QUEUE"), req.Message)
	if err != nil {
		return SendResponse{
			Sent: false,
		}, err
	}

	return SendResponse{
		Sent: true,
	}, nil

}
