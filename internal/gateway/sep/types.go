package sep

type Gate struct {
	MerchantID string
	TerminalID string
}

type tokenRequest struct {
	TerminalID  string `json:"TerminalId"`
	Action      string `json:"Action"`
	RedirectURL string `json:"RedirectURL"`
	CellNumber  string `json:"CellNumber"`
	Amount      uint   `json:"Amount"`
	ResNum      uint   `json:"ResNum"`
}

type tokenResp struct {
	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDesc"`
	Token     string `json:"token"`
	Status    int    `json:"status"`
}

type ConfirmResp struct {
	Amount           string `schema:"Amount"`
	HashedCardNumber string `schema:"HashedCardNumber"`
	Mid              string `schema:"MID"`
	RefNum           string `schema:"RefNum"`
	PaymentID        string `schema:"ResNum"`
	RRN              string `schema:"Rrn"`
	SecurePan        string `schema:"SecurePan"`
	State            string `schema:"State"`
	Status           string `schema:"Status"`
	TerminalId       string `schema:"TerminalId"`
	Token            string `schema:"Token"`
	TraceNo          string `schema:"TraceNo"`
}

type VerifyResp struct {
	Body struct {
		VerifyTransactionResponse struct {
			Result int `xml:"result"`
		} `xml:"verifyTransactionResponse"`
	} `xml:"Body"`
	Description string
}
