package atipay

type Gate struct {
	ApiKey     string
	TerminalID string
}

type tokenRequest struct {
	ApiKey      string `json:"apiKey"`
	RedirectURL string `json:"redirectUrl"`
	CellNumber  string `json:"cellNumber"`
	Description string `json:"description"`
	Amount      uint   `json:"amount"`
	ResNum      uint   `json:"invoiceNumber"`
}

type tokenResp struct {
	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDescription"`
	Token     string `json:"token"`
	Status    string `json:"status"`
}

type VerifyResp struct {
	Amount float32 `json:"amount"`
}

type VerifyRequest struct {
	TerminalID      string `json:"terminalId"`
	ReferenceNumber string `json:"referenceNumber"`
	ApiKey          string `json:"apiKey"`
}
