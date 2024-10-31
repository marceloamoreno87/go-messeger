package domain

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

/*
Estrutura Message representa uma mensagem a ser enviada.
Campos:
- JID: Identificador do remetente.
- To: Identificador do destinatário.
- Message: Conteúdo da mensagem a ser enviada.
*/
type Message struct {
	JID     string `json:"jid"`
	To      string `json:"to"`
	Message string `json:"message"`
}

/*
Método GetMessage retorna o conteúdo da mensagem.
Retorna:
- Uma string contendo o conteúdo da mensagem.
*/
func (m *Message) GetMessage() string {
	return m.Message
}

/*
Estrutura SendMessage contém o repositório WhatsAppRepository.
Esta estrutura é responsável por enviar mensagens usando o serviço WhatsApp.
*/
type SendMessage struct {
	WhatsAppRepository WhatsAppRepository
}

/*
Método Send envia uma mensagem usando o serviço WhatsApp.
Parâmetros:
- message: Ponteiro para a estrutura Message contendo os detalhes da mensagem a ser enviada.
Retorna:
- Um erro, se houver.
*/
func (s *SendMessage) Send(message *Message) (err error) {
	/*
	   Converte o identificador do remetente (JID) para o formato types.JID.
	*/
	JID := NumberToJID(message.JID)

	/*
	   Encontra o dispositivo WhatsApp associado ao JID.
	   Se o dispositivo não for encontrado, retorna um erro.
	*/
	deviceStore, err := s.WhatsAppRepository.FindDeviceWM(context.Background(), JID)
	if err != nil {
		return err
	}

	/*
	   Cria um novo cliente WhatsApp usando o dispositivo encontrado.
	   Conecta ao cliente WhatsApp.
	   Se a conexão falhar, retorna um erro.
	*/
	client := whatsmeow.NewClient(deviceStore, nil)
	if err = client.Connect(); err != nil {
		return err
	}

	/*
	   Constrói o identificador do destinatário (TO) no formato types.JID.
	*/
	TO := types.JID{
		Server: "s.whatsapp.net",
		User:   message.To,
	}

	/*
	   Envia a mensagem para o destinatário usando o cliente WhatsApp.
	   Se o envio falhar, retorna um erro.
	*/
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
