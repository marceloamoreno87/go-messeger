package domain

type ValidateRequest struct {
	JID string `json:"phone_number" validate:"required"`
}

type ValidateResponse struct {
	Active bool `json:"active"`
}

type SendRequest struct {
	JID          string `json:"jid"`
	To           string `json:"to"`
	Title        string `json:"title"`
	Service      string `json:"service"`
	When         string `json:"when"`
	Duration     string `json:"duration"`
	Professional string `json:"professional"`
	Code         string `json:"code"`
	Footer       string `json:"footer"`
}

type SendResponse struct {
	Sent bool `json:"sent"`
}
