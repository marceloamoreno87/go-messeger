package domain

import (
	"encoding/json"
	"net/http"
)

type WhatsAppHandler struct {
	WhatsAppService WhatsAppService
}

func (h WhatsAppHandler) Connect(w http.ResponseWriter, r *http.Request) {

	data, err := h.WhatsAppService.Connect(r.Context())
	if err != nil && data == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=qr.png")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (h WhatsAppHandler) Validate(w http.ResponseWriter, r *http.Request) {

	req := ValidateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.WhatsAppService.Validate(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)

}

func (h WhatsAppHandler) Send(w http.ResponseWriter, r *http.Request) {
	req := SendRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.WhatsAppService.Send(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
