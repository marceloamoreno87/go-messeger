package domain

import (
	"context"
	"encoding/json"
	"net/http"
)

/*
Estrutura WhatsAppHandler que contém o serviço WhatsAppService.
Esta estrutura é responsável por lidar com as solicitações HTTP relacionadas ao WhatsApp.
*/
type WhatsAppHandler struct {
	WhatsAppService WhatsAppService
}

/*
Método Connect lida com a solicitação HTTP para conectar ao serviço WhatsApp.
Gera um código QR para autenticação e o retorna como uma imagem PNG.
Em caso de erro, retorna um status HTTP 500.
*/
func (h WhatsAppHandler) Connect(w http.ResponseWriter, r *http.Request) {
	res, err := h.WhatsAppService.Connect(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

/*
Método Validate lida com a solicitação HTTP para validar um número de telefone.
Decodifica a solicitação JSON para a estrutura ValidateRequest.
Em caso de erro, retorna um status HTTP 400.
Chama o serviço de validação e retorna a resposta como JSON.
Em caso de erro no serviço, retorna um status HTTP 500.
*/
func (h WhatsAppHandler) Validate(w http.ResponseWriter, r *http.Request) {
	req := ValidateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.WhatsAppService.Validate(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

/*
Método Send lida com a solicitação HTTP para enviar uma mensagem.
Decodifica a solicitação JSON para a estrutura SendRequest.
Em caso de erro, retorna um status HTTP 400.
Chama o serviço de envio de mensagem e retorna a resposta como JSON.
Em caso de erro no serviço, retorna um status HTTP 500.
*/
func (h WhatsAppHandler) Send(w http.ResponseWriter, r *http.Request) {
	req := SendRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.WhatsAppService.Send(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
