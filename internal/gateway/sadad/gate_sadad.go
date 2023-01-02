package sadad

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/conv"
	"gitlab.boomerangapp.ir/back/utils/pkg/io"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

func (*Gate) MakeRedirectURL(token string) string {

	return consts.SadadPSPURL + "?Token=" + token
}

func (g *Gate) GetToken(pr types.PaymentRequest, mobile string) (string, error) {

	key, err := base64.StdEncoding.DecodeString(g.ApiKey)
	if err != nil {
		return "", err
	}

	signData, err := tripleEcbDesEncrypt(conv.S2B(g.TerminalID+";"+
		strconv.Itoa(int(pr.ID))+";"+
		strconv.Itoa(int(pr.Amount))), key)
	if err != nil {
		return "", err
	}

	request := tokenRequest{
		MerchantID:    g.MerchantID,
		TerminalID:    g.TerminalID,
		Amount:        pr.Amount,
		OrderID:       pr.ID,
		Phone:         mobile,
		LocalDateTime: time.Now().Format("01/02/2006 15:04:05 PM"),
		ReturnURL: configs.Get().AppConfigs.BaseURL + "/" + consts.HttpPrefix +
			consts.ReturnURL + "/" + strconv.Itoa(int(pr.ID)),
		SignData: base64.StdEncoding.EncodeToString(signData),
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(consts.SadadTokenURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer io.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("خطا در ارتباط با سداد: " + resp.Status)
	}

	var result tokenResp
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", errors.New("خطا در بازخوانی اطلاعات سداد" + err.Error())
	}

	if result.Token != "" {
		return result.Token, nil
	}

	return "", errors.New(result.Description + "---" + result.ResCode)
}

func (g *Gate) Verify(token string, resp interface{}) error {

	request := VerifyRequest{
		Token: token,
	}

	k, err := base64.StdEncoding.DecodeString(g.ApiKey)
	if err != nil {
		return err
	}
	d, err := tripleEcbDesEncrypt([]byte(token), k)
	if err != nil {
		return err
	}

	request.SignData = base64.StdEncoding.EncodeToString(d)
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	var response *http.Response
	for {
		response, err = http.Post(consts.SadadVerifyURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			sentrylog.Warning(errors.New("خطا در ارتباط با سداد"))
			continue
		}
		if response.StatusCode != http.StatusOK {
			sentrylog.Warning(errors.New("خطا در ارتباط با سداد" + response.Status))
			continue
		}
		break
	}
	defer io.SafeClose(response.Body)

	var svr VerifyResp
	err = json.NewDecoder(response.Body).Decode(&svr)
	if err != nil {
		return errors.New("خطا در بازخوانی اطلاعات سداد")
	}

	reflect.ValueOf(resp).Elem().Set(reflect.ValueOf(svr))
	return nil
}

func (*Gate) CheckConfirmResp(resp map[string]string) (string, bool) {

	return "", false
}

func (*Gate) CheckVerifyResp(resp interface{}) bool {

	v := resp.(VerifyResp)
	return v.ResCode == 0
}

func (*Gate) SerializeToDB(resp interface{}) string {

	var result string
	switch v := resp.(type) {
	case ConfirmResp:
		result = fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s",
			"IsWalletPayment="+v.IsWalletPayment, "PrimaryAccNo="+v.PrimaryAccNo,
			"HashedCardNo="+v.HashedCardNo, "__RequestVerificationToken="+v.RequestVerificationToken,
			"SwitchResCode="+v.SwitchResCode, "OrderId="+v.PaymentID,
			"ResCode="+v.ResCode, "token="+v.Token)

	case VerifyResp:
		result = fmt.Sprintf("%s;%s;%s;%s;%s;%s",
			"ResCode="+string(rune(v.ResCode)), "SystemTraceNo="+v.SystemTraceNo,
			"RetrivalRefNo="+v.RetrivalRefNo, "Amount="+strconv.Itoa(int(v.Amount)),
			"Description="+v.Description, "OrderId="+strconv.Itoa(int(v.PaymentID)))
	}

	return result
}
