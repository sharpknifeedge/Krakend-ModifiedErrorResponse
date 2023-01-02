package consts

import "gitlab.boomerangapp.ir/back/utils/pkg/response"

const ActionType = "increase_wallet"

const (
	StatusIngress uint8 = iota
	StatusVerified
	StatusFailed
)

//sadad consts
const (
	SadadTokenURL  = "https://sadad.shaparak.ir/api/v0/Request/PaymentRequest"
	SadadVerifyURL = "https://sadad.shaparak.ir/api/v0/Advice/Verify"
	SadadPSPURL    = "https://sadad.shaparak.ir/Purchase"
)

//sep consts
const (
	SepTokenURL  = "https://sep.shaparak.ir/MobilePG/MobilePayment"
	SepVerifyURL = "https://sep.shaparak.ir/payments/referencepayment.asmx"
	SepPSPForm   = `<html lang="fa">
<head>
   <meta charset="UTF-8">
   <title>بومرنگ</title>
</head>
<body onload="document.forms['form'].submit()">
<p>در حال انتقال به درگاه سپ...</p>
<form name="form" action="https://sep.shaparak.ir/OnlinePG/OnlinePG" method="post">
   <input name="Token" type="hidden" value="%s" />
   <input name="GetMethod" type="hidden" value="false">
</form>
</body>
</html>`
)

//atipay consts
const (
	AtipayTokenURL  = "https://mipg.atipay.net/v1/get-token"
	AtipayVerifyURL = "https://mipg.atipay.net/v1/verify-payment"
	AtipayPSPForm   = `<!Doctype html>
<html lang="fa">
<head>
   <meta charset="UTF-8">
   <title>بومرنگ</title>
</head>
<body onload="document.forms['form'].submit()">
<p>در حال انتقال به درگاه آتی پی...</p>
<form name="form" action="https://mipg.atipay.net/v1/redirect-to-gateway" method="post">
   <input type="hidden" value="%s" name="token">
</form>
</body>
</html>`
)

//pec consts
const (
	PecTokenURL  = "https://pec.shaparak.ir/NewIPGServices/Sale/SaleService.asmx"
	PecVerifyURL = "https://pec.shaparak.ir/NewIPGServices/Confirm/ConfirmService.asmx"
	PecPSPURL    = "https://pec.shaparak.ir/NewIPG/?token="
)

const (
	//All the psp gateways
	PSPSadad         = "sadad"
	PSPSep           = "sep"
	PSPAtipay        = "atipay"
	PSPPec           = "pec"
	PSPPecPeraion    = "بانک پارسیان"
	PSPAtipayPeraion = "آتی پی"
	PSPSadadPersion  = "بانک ملی"
	PSPSepPersion    = "بانک سامان"
)

//routes
const (
	HttpPrefix    = "/payment"
	RouteRedirect = "/redirectToPSP"
	ReturnURL     = "/confirm"
)

//service error codes
const (
	ErrCodeInvalidPSP uint16 = iota + 900
	ErrCodeInvalidTerminal
	ErrCodeUnableToGetToken
	ErrCodeInvalidUserID
	ErrCodeUserIDConflict
	ErrCodeLessThanMinAmount
	ErrCodeMoreThanMaxBalance
	ErrCodeMoreThanDailyDeposit
)

var Messages = map[uint16]response.Message{
	ErrCodeInvalidPSP:           {Header: "درگاه نامعتبر است"},
	ErrCodeInvalidTerminal:      {Header: "شماره درگاه نامعتبر است"},
	ErrCodeUnableToGetToken:     {Header: "خطا در دریافت توکن از درگاه"},
	ErrCodeInvalidUserID:        {Header: "خطا در خواندن اطلاعات کاربر"},
	ErrCodeUserIDConflict:       {Header: "عدم تطابق آیدی"},
	ErrCodeLessThanMinAmount:    {Header: "مبلغ کمتر از حد مجاز"},
	ErrCodeMoreThanMaxBalance:   {Header: "تراکنش خارج از سقف موجودی کیف پول"},
	ErrCodeMoreThanDailyDeposit: {Header: "تراکنش خارج از سقف پرداخت روزانه"},
}
