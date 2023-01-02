package pec

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/io"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

func (*Gate) MakeRedirectURL(token string) string {

	return consts.PecPSPURL + token
}

func (g *Gate) GetToken(pr types.PaymentRequest, mobile string) (string, error) {

	body := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
    <Body>
        <SalePaymentRequest xmlns="https://pec.Shaparak.ir/NewIPGServices/Sale/SaleService">
            <requestData>
                <LoginAccount>` + g.ApiKey + `</LoginAccount>
                <Amount>` + strconv.Itoa(int(pr.Amount)) + `</Amount>
                <OrderId>` + strconv.Itoa(int(pr.ID)) + `</OrderId>
				<CallBackUrl>` + configs.Get().AppConfigs.BaseURL + consts.HttpPrefix +
		consts.ReturnURL + "/" + strconv.Itoa(int(pr.ID)) + `</CallBackUrl>
				<Originator>` + strings.TrimPrefix(mobile, "0") + `</Originator>
            </requestData>
        </SalePaymentRequest>
    </Body>
</Envelope>`

	resp, err := http.Post(consts.PecTokenURL, "text/xml", bytes.NewReader([]byte(body)))
	if err != nil {
		return "", err
	}
	defer io.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("خطا در ارتباط با پارسیان : " + resp.Status)
	}

	var result tokenResp
	err = xml.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", errors.New("خطا در بازخوانی اطلاعات پارسیان " + err.Error())
	}

	if len(result.Body.Resp.Result.Token) > 0 && result.Body.Resp.Result.Status == 0 {
		return result.Body.Resp.Result.Token, nil
	}

	return "", errors.New(result.Body.Resp.Result.Message + "***" + strconv.Itoa(result.Body.Resp.Result.Status))
}

func (*Gate) SerializeToDB(resp interface{}) string {

	var result string
	switch v := resp.(type) {
	case map[string]string:
		result = strings.Join([]string{
			"Amount=" + v["Amount"],
			"DiscoutedProduct=" + v["DiscoutedProduct"],
			"HashCardNumber=" + v["HashCardNumber"],
			"OrderId=" + v["OrderId"],
			"RRN=" + v["RRN"],
			"STraceNo=" + v["STraceNo"],
			"SwAmount=" + v["SwAmount"],
			"TerminalNo=" + v["TerminalNo"],
			"Token=" + v["Token"],
			"TspToken=" + v["TspToken"],
			"status=" + v["status"],
		}, ";")
	case VerifyResp:
		result = strings.Join([]string{
			"Status=" + strconv.Itoa(v.Body.Resp.Result.Status),
			"CardNumberMasked=" + v.Body.Resp.Result.CardNumberMasked,
			"Token=" + strconv.Itoa(int(v.Body.Resp.Result.Token)),
			"RRN=" + strconv.Itoa(int(v.Body.Resp.Result.RRN)),
		}, ";")
	}

	return result
}

func (*Gate) CheckConfirmResp(resp map[string]string) (string, bool) {

	if resp["status"] != "0" || resp["RRN"] == "0" {
		return "", false
	}

	return resp["Token"], true
}

func (*Gate) CheckVerifyResp(resp interface{}) bool {

	v := resp.(VerifyResp)
	if v.Body.Resp.Result.Status == 0 || v.Body.Resp.Result.Status == -1533 {
		return true
	}

	return false
}

func (g *Gate) Verify(refNum string, resp interface{}) error {

	body := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
    <Body>
        <ConfirmPayment xmlns="https://pec.Shaparak.ir/NewIPGServices/Confirm/ConfirmService">
            <requestData>
                <LoginAccount>` + g.ApiKey + `</LoginAccount>
                <Token>` + refNum + `</Token>
            </requestData>
        </ConfirmPayment>
    </Body>
</Envelope>`

	var response *http.Response
	var err error
	try := 0
	for ; try < 3; try++ {
		response, err = http.Post(consts.PecVerifyURL, "text/xml", bytes.NewReader([]byte(body)))
		if err != nil {
			sentrylog.Warning(errors.New("خطا در ارتباط با پارسیان " + err.Error()))
			continue
		}
		if response.StatusCode != http.StatusOK {
			sentrylog.Warning(errors.New("خطا در ارتباط با پارسیان " + response.Status))
			continue
		}
		break
	}
	defer io.SafeClose(response.Body)

	if try == 2 && (response == nil || response.StatusCode != http.StatusOK) {
		return errors.New("خطا در ارتباط با پارسیان " + response.Status)
	}

	var svr VerifyResp
	err = xml.NewDecoder(response.Body).Decode(&svr)
	if err != nil {
		return errors.New("خطا در بازخوانی اطلاعات پارسیان " + err.Error())
	}

	reflect.ValueOf(resp).Elem().Set(reflect.ValueOf(svr))
	return nil
}
