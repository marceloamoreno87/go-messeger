package domain

type ValidateRequest struct {
	JID string `json:"phone_number" validate:"required"`
}

type ValidateResponse struct {
	Active bool `json:"active"`
}

type SendRequest struct {
	Message string `json:"message" validate:"required"`
}

type SendResponse struct {
	Sent bool `json:"sent"`
}
