package sep

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/conv"
	"gitlab.boomerangapp.ir/back/utils/pkg/io"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

func GetForm(token string) []byte {

	form := fmt.Sprintf(consts.SepPSPForm, token)
	return conv.S2B(form)
}

func (*Gate) MakeRedirectURL(token string) string {

	return configs.Get().AppConfigs.BaseURL + consts.HttpPrefix +
		consts.RouteRedirect + "/" + consts.PSPSep + "/" + token
}

func (g *Gate) GetToken(pr types.PaymentRequest, mobile string) (string, error) {

	request := tokenRequest{
		TerminalID: g.TerminalID,
		Action:     "token",
		Amount:     pr.Amount,
		ResNum:     pr.ID,
		CellNumber: mobile,
		RedirectURL: configs.Get().AppConfigs.BaseURL + "/" + consts.HttpPrefix +
			consts.ReturnURL + "/" + strconv.Itoa(int(pr.ID)),
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(consts.SepTokenURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer io.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("خطا در ارتباط با سپ: " + resp.Status)
	}

	var result tokenResp
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", errors.New("خطا در بازخوانی اطلاعات سپ" + err.Error())
	}

	if result.Token != "" {
		return result.Token, nil
	}

	return "", errors.New(result.ErrorDesc + "---" + result.ErrorCode)
}

func (*Gate) SerializeToDB(resp interface{}) string {

	var result string
	switch v := resp.(type) {
	case ConfirmResp:
		result = fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s",
			"Amount="+v.Amount, "HashedCardNumber="+v.HashedCardNumber,
			"MID="+v.Mid, "RefNum="+v.RefNum,
			"ResNum="+v.PaymentID, "Rrn="+v.RRN,
			"SecurePan="+v.SecurePan, "State="+v.State,
			"Status="+v.Status, "TerminalId="+v.TerminalId,
			"Token="+v.Token, "TraceNo="+v.TraceNo)
	case VerifyResp:
		result = fmt.Sprintf("result=%d;description=%s",
			v.Body.VerifyTransactionResponse.Result,
			v.Description)
	}

	return result
}

func (*Gate) CheckConfirmResp(resp map[string]string) (string, bool) {

	return "", false
}

func (*Gate) CheckVerifyResp(resp interface{}) bool {

	v := resp.(VerifyResp)
	return v.Body.VerifyTransactionResponse.Result > 0
}

func (g *Gate) Verify(refNum string, resp interface{}) error {

	body := `<soapenv:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:Foo">
<soapenv:Header/>
<soapenv:Body>
   <urn:verifyTransaction soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
      <String_1 xsi:type="xsd:string">` + refNum + `</String_1>
      <String_2 xsi:type="xsd:string">` + g.MerchantID + `</String_2>
   </urn:verifyTransaction>
</soapenv:Body>
</soapenv:Envelope>`

	var response *http.Response
	for {
		var err error
		response, err = http.Post(consts.SepVerifyURL, "text/xml", bytes.NewBufferString(body))
		if err != nil {
			sentrylog.Warning(errors.New("خطا در ارتباط با سپ"))
			continue
		}
		if response.StatusCode != http.StatusOK {
			sentrylog.Warning(errors.New("خطا در ارتباط با سپ" + response.Status))
			continue
		}
		break
	}
	defer io.SafeClose(response.Body)

	var svr VerifyResp
	err := xml.NewDecoder(response.Body).Decode(&svr)
	if err != nil {
		return errors.New("خطا در پردازش پاسخ سپ")
	}

	switch svr.Body.VerifyTransactionResponse.Result {
	case -1:
		svr.Description = "خطا در پردازش اطلاعات ارسالی"
	case -3:
		svr.Description = "ورودی ها حاوی کاراکتر غیر مجاز است"
	case -4:
		svr.Description = "کلمه‌ی عبور و یا کد فروشنده اشتباه است"
	case -6:
		svr.Description = "سند برگشتی و یا منقضی است"
	case -7:
		svr.Description = "رسید دیجیتال تهی است"
	case -8:
		svr.Description = "طول ورودی ها بیشتر از حد مجاز است"
	case -9:
		svr.Description = "وجود کاراکترهای غیر مجاز در مبلغ برگشتی"
	case -10:
		svr.Description = "رسید دیجیتال حاوی کاراکتر غیر مجاز است"
	case -11:
		svr.Description = "طول ورودی ها کمتر ازحد مجاز است"
	case -12:
		svr.Description = "مبلغ برگشتی منفی است"
	case -13:
		svr.Description = "مبلغ برگشت جزئی بیش از حد مجاز است"
	case -14:
		svr.Description = "تراکنش تعریف نشده است"
	case -15:
		svr.Description = "مبلغ برگشتی به صورت اعشاری است"
	case -16:
		svr.Description = "خطای داخلی سیستم"
	case -17:
		svr.Description = "برگشت زدن جزيی تراکنش مجاز نمیباشد"
	case -18:
		svr.Description = "IP کارگزار تعریف نشده است"
	}

	reflect.ValueOf(resp).Elem().Set(reflect.ValueOf(svr))
	return nil
}
