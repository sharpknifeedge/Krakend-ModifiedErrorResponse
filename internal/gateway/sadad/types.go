package sadad

type Gate struct {
	MerchantID string
	TerminalID string
	ApiKey     string
}

type tokenRequest struct {
	MerchantID    string `json:"MerchantId"`
	TerminalID    string `json:"TerminalId"`
	LocalDateTime string `json:"LocalDateTime"`
	ReturnURL     string `json:"ReturnUrl"`
	SignData      string `json:"SignData"`
	Phone         string `json:"UserId"`
	Amount        uint   `json:"Amount"`
	OrderID       uint   `json:"OrderId"`
}

type tokenResp struct {
	ResCode     string `json:"ResCode"`
	Token       string `json:"Token"`
	Description string `json:"Description"`
}

type ConfirmResp struct {
	PaymentID                string `schema:"OrderId"`
	HashedCardNo             string `schema:"HashedCardNo"`
	PrimaryAccNo             string `schema:"PrimaryAccNo"`
	SwitchResCode            string `schema:"SwitchResCode"`
	ResCode                  string `schema:"ResCode"`
	Token                    string `schema:"token"`
	RequestVerificationToken string `schema:"__RequestVerificationToken"`
	IsWalletPayment          string `schema:"IsWalletPayment"`
}

type VerifyResp struct {
	Description   string `json:"Description"`
	RetrivalRefNo string `json:"RetrivalRefNo"`
	SystemTraceNo string `json:"SystemTraceNo"`
	Amount        uint   `json:"Amount"`
	PaymentID     uint   `json:"OrderId"`
	ResCode       int    `json:"ResCode"`
}

type VerifyRequest struct {
	Token    string `json:"Token"`
	SignData string `json:"SignData"`
}
