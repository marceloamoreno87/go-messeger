package domain

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gonext/core"
	"log"
	"os"
	"strconv"

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
func (s WhatsAppService) Connect(ctx context.Context) (res ConnectResponse, err error) {

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
			//request para atualizar o agenda que o aparelho foi configurado com sucesso

		}
	}()

	store.SetOSInfo("Windows", [3]uint32{1, 2, 3})

	qrChan, _ := client.GetQRChannel(context.Background())
	err = client.Connect()
	if err != nil {
		return ConnectResponse{}, err
	}
	qr := []byte{}
	for evt := range qrChan {
		if evt.Event == "code" {
			qr, err = qrcode.Encode(evt.Code, qrcode.Medium, 256)
			if err != nil {
				return ConnectResponse{}, err
			}
			break
		}
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qr)
	qrToStringBase64 := fmt.Sprintf("data:image/png;base64,%s", qrBase64)
	return ConnectResponse{
		AuthCode:  qrToStringBase64,
		SessionID: strconv.FormatUint(uint64(deviceStore.RegistrationID), 10),
	}, nil
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
	sessionId := req.SessionId
	_, err = s.WhatsAppRepository.FindDeviceWM(ctx, sessionId)
	if err != nil && err != sql.ErrNoRows {
		return ValidateResponse{
			Active: false,
		}, err
	}

	if err == sql.ErrNoRows {
		return ValidateResponse{
			Active: false,
		}, nil
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

func (s WhatsAppService) CreateAccount(ctx context.Context, req CreateAccountRequest) (res CreateAccountResponse, err error) {
	name := req.Name
	origin := req.Origin
	externalId := req.ExternalId

	id, err := s.WhatsAppRepository.CreateAccount(ctx, name, origin, externalId)
	if err != nil {
		return CreateAccountResponse{}, err
	}

	return CreateAccountResponse{
		ID: id,
	}, nil
}

func (s WhatsAppService) CreateSession(ctx context.Context, req CreateSessionRequest) (res CreateSessionResponse, err error) {

	deviceStore := s.WhatsAppRepository.CreateDeviceWM(context.Background())
	client := whatsmeow.NewClient(deviceStore, nil)

	pairSuccessChan := make(chan *events.PairSuccess, 1)
	client.AddEventHandler(func(evt interface{}) {
		if pairSuccess, ok := evt.(*events.PairSuccess); ok {
			pairSuccessChan <- pairSuccess
		}
	})

	pairErrorChan := make(chan *events.PairError, 1)
	client.AddEventHandler(func(evt interface{}) {
		if PairError, ok := evt.(*events.PairError); ok {
			pairErrorChan <- PairError
		}
	})

	sessionId := strconv.FormatUint(uint64(deviceStore.RegistrationID), 10)

	go func() {
		select {
		case <-pairSuccessChan:
			//Atualiza o sessioin para CONNECTED
			err := s.WhatsAppRepository.UpdateSession(ctx, sessionId, "CONNECTED", deviceStore.ID.String())
			if err != nil {
				log.Fatalf("Failed to update session: %v", err)
			}

		case <-pairErrorChan:
			// apaga o accountIdv2 e sessionIdv2 no workspace
			// apaga o account
			// apaga o session
		}
	}()

	store.SetOSInfo("Windows", [3]uint32{1, 2, 3})

	qrChan, _ := client.GetQRChannel(context.Background())
	err = client.Connect()
	if err != nil {
		return CreateSessionResponse{}, err
	}
	qr := []byte{}
	for evt := range qrChan {
		if evt.Event == "code" {
			qr, err = qrcode.Encode(evt.Code, qrcode.Medium, 256)
			if err != nil {
				return CreateSessionResponse{}, err
			}
			break
		}
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qr)
	qrToStringBase64 := fmt.Sprintf("data:image/png;base64,%s", qrBase64)

	accountID := req.AccountID

	_, _, err = s.WhatsAppRepository.CreateSession(ctx, sessionId, accountID, qrToStringBase64)
	if err != nil {
		return CreateSessionResponse{}, err
	}

	return CreateSessionResponse{
		AuthCode:  qrToStringBase64,
		SessionID: sessionId,
	}, nil

}

func (s WhatsAppService) GetSessionByID(ctx context.Context, id string) (res GetSessionByIDResponse, err error) {
	res, err = s.WhatsAppRepository.GetSessionByID(ctx, id)
	if err != nil {
		return GetSessionByIDResponse{}, err
	}

	return res, nil
}

func (s WhatsAppService) DeleteSession(ctx context.Context, id string) (err error) {

	session, err := s.WhatsAppRepository.GetSessionByID(ctx, id)
	if err != nil {
		return err
	}

	deviceStore, err := s.WhatsAppRepository.FindDeviceWM(context.Background(), session.ID)
	if err != nil {
		return err
	}

	err = deviceStore.Delete()
	if err != nil {
		return err
	}

	err = s.WhatsAppRepository.DeleteSession(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s WhatsAppService) DeleteAccount(ctx context.Context, id string) (err error) {
	err = s.WhatsAppRepository.DeleteAccount(ctx, id)
	if err != nil {
		return err

	}
	return nil
}
