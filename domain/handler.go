package domain

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (h WhatsAppHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	req := CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.WhatsAppService.CreateAccount(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (h WhatsAppHandler) CreateSession(w http.ResponseWriter, r *http.Request) {

	req := CreateSessionRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.WhatsAppService.CreateSession(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (h WhatsAppHandler) GetSessionByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	res, err := h.WhatsAppService.GetSessionByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (h WhatsAppHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	err := h.WhatsAppService.DeleteSession(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h WhatsAppHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	err := h.WhatsAppService.DeleteAccount(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
