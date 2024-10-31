package domain

/*
Estrutura ValidateRequest representa a solicitação para validar um número de telefone.
Campos:
- JID: Número de telefone a ser validado. Este campo é obrigatório.
*/
type ValidateRequest struct {
	JID string `json:"phone_number" validate:"required"`
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
	JID     string `json:"jid"`
	To      string `json:"to"`
	Message string `json:"message"`
}

/*
	Estrutura SendResponse representa a resposta para a solicitação de envio de mensagem.
	Campos:
	- Sent: Indica se a mensagem foi enviada com sucesso.
*/
type SendResponse struct {
	Sent bool `json:"sent"`
}
