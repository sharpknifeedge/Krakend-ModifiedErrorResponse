package pec

type Gate struct {
	ApiKey string
}

type tokenResp struct {
	Body struct {
		Resp struct {
			Result struct {
				Token   string `xml:"Token"`
				Message string `xml:"Message"`
				Status  int    `xml:"Status"`
			} `xml:"SalePaymentRequestResult"`
		} `xml:"SalePaymentRequestResponse"`
	} `xml:"Body"`
}

type VerifyResp struct {
	Body struct {
		Resp struct {
			Result struct {
				CardNumberMasked string `xml:"CardNumberMasked"`
				Status           int    `xml:"Status"`
				Token            uint   `xml:"Token"`
				RRN              uint   `xml:"RRN"`
			} `xml:"ConfirmPaymentResult"`
		} `xml:"ConfirmPaymentResponse"`
	} `xml:"Body"`
}
