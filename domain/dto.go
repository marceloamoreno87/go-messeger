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

type CreateAccountRequest struct {
	Name       string `json:"name"`
	Origin     string `json:"origin"`
	ExternalId string `json:"externalId"`
}

type CreateAccountResponse struct {
	ID string `json:"id"`
}

type CreateSessionRequest struct {
	AccountID string `json:"accountId"`
}

type CreateSessionResponse struct {
	AuthCode  string `json:"authCode"`
	SessionID string `json:"sessionId"`
}

type SessionInfo struct {
	Name            *string `json:"name,omitempty"`
	PhoneID         *string `json:"phoneId,omitempty"`
	PhoneSerialized *string `json:"phoneSerialized,omitempty"`
	Platform        *string `json:"platform,omitempty"`
}

type GetSessionByIDResponse struct {
	ID            string       `json:"id"`
	AccountID     string       `json:"accountId"`
	SessionInfo   *SessionInfo `json:"sessionInfo"`
	Status        string       `json:"status"`
	AuthCode      string       `json:"authCode"`
	ReadyAt       *string      `json:"readyAt,omitempty"`
	FailureReason *string      `json:"failureReason,omitempty"`
	CreatedAt     string       `json:"createdAt"`
	UpdatedAt     string       `json:"updatedAt"`
}
