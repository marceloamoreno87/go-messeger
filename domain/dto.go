package domain

type ValidateRequest struct {
	JID string `json:"phone_number" validate:"required"`
}

type ValidateResponse struct {
	Active bool `json:"active"`
}

type SendRequest struct {
	JID     string `json:"jid"`
	To      string `json:"to"`
	Message string `json:"message"`
}

type SendResponse struct {
	Sent bool `json:"sent"`
}
