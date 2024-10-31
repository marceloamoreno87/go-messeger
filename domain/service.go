package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"gonext/core"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types/events"
)

/*
Estrutura WhatsAppService que contém o repositório WhatsAppRepository e o serviço RabbitMQ.
Esta estrutura é responsável por fornecer funcionalidades relacionadas ao WhatsApp.
*/
type WhatsAppService struct {
	WhatsAppRepository WhatsAppRepository
	RabbitMQService    *core.RabbitMQClient
}

/*
Método Connect lida com a conexão ao serviço WhatsApp.
Gera um código QR para autenticação e o retorna como uma imagem PNG.
Em caso de sucesso, atualiza ou cria o dispositivo no banco de dados Postgres.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
Retorna:
- Um slice de bytes contendo a imagem PNG do código QR e um erro, se houver.
*/
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
			JID := fmt.Sprintf("%s:%d@%s", deviceStore.ID.User, deviceStore.ID.Device, deviceStore.ID.Server)
			err := s.WhatsAppRepository.UpdateOrCreatePG(context.Background(), JID, deviceStore.ID.User)
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

/*
Método Validate lida com a validação de um número de telefone.
Verifica se o dispositivo associado ao número de telefone está ativo.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
- req: Estrutura ValidateRequest contendo o número de telefone a ser validado.
Retorna:
- Uma estrutura ValidateResponse indicando se o número de telefone está ativo e um erro, se houver.
*/
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

/*
Método Send lida com o envio de uma mensagem.
Publica a mensagem na fila RabbitMQ.
Parâmetros:
- ctx: Contexto para controle de cancelamento e prazos.
- req: Estrutura SendRequest contendo os detalhes da mensagem a ser enviada.
Retorna:
- Uma estrutura SendResponse indicando se a mensagem foi enviada com sucesso e um erro, se houver.
*/
func (s WhatsAppService) Send(ctx context.Context, req SendRequest) (res SendResponse, err error) {
	client, err := s.RabbitMQService.Connect()
	if err != nil {
		return SendResponse{
			Sent: false,
		}, err
	}
	defer client.Close()

	jsonReq, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	_, err = client.Publish(os.Getenv("RABBITMQ_QUEUE"), string(jsonReq))
	if err != nil {
		return SendResponse{
			Sent: false,
		}, err
	}

	return SendResponse{
		Sent: true,
	}, nil
}
