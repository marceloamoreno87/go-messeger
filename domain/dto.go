package domain

/*
Estrutura ValidateRequest representa a solicitação para validar um número de telefone.
Campos:
- JID: Número de telefone a ser validado. Este campo é obrigatório.
*/
type ValidateRequest struct {
	SessionId string `json:"sessionId" validate:"required"`
}

/*
Estrutura ValidateResponse representa a resposta para a solicitação de validação.
Campos:
- Active: Indica se o número de telefone está ativo.
*/
type ValidateResponse struct {
	Active bool `json:"active"`
}

/*
Estrutura SendRequest representa a solicitação para enviar uma mensagem.
Campos:
- JID: Identificador do remetente.
- To: Identificador do destinatário.
- Message: Conteúdo da mensagem a ser enviada.
*/
type SendRequest struct {
	SessionId string `json:"sessionId"`
	To        string `json:"to"`
	Message   string `json:"message"`
}

/*
Estrutura SendResponse representa a resposta para a solicitação de envio de mensagem.
Campos:
- Sent: Indica se a mensagem foi enviada com sucesso.
*/
type SendResponse struct {
	Sent bool `json:"sent"`
}

/*
Estrutura ConnectRequest representa a solicitação para conectar ao serviço.
Campos:
- JID: Número de telefone a ser conectado.
*/
type ConnectResponse struct {
	AuthCode  string `json:"authCode"`
	SessionID string `json:"sessionId"`
}
