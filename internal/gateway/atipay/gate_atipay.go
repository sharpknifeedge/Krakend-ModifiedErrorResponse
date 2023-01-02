package atipay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/conv"
	"gitlab.boomerangapp.ir/back/utils/pkg/io"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

func GetForm(token string) []byte {

	return conv.S2B(fmt.Sprintf(consts.AtipayPSPForm, token))
}

func (*Gate) MakeRedirectURL(token string) string {

	return configs.Get().AppConfigs.BaseURL + consts.HttpPrefix +
		consts.RouteRedirect + "/" + consts.PSPAtipay + "/" + token
}

func (g *Gate) GetToken(pr types.PaymentRequest, mobile string) (string, error) {

	request := tokenRequest{
		ApiKey:      g.ApiKey,
		Amount:      pr.Amount,
		ResNum:      pr.ID,
		CellNumber:  mobile,
		Description: "کار بر به شماره موبایل " + mobile,
		RedirectURL: configs.Get().AppConfigs.BaseURL + consts.HttpPrefix +
			consts.ReturnURL + "/" + strconv.Itoa(int(pr.ID)),
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(consts.AtipayTokenURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer io.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("خطا در ارتباط با آتی پی: " + resp.Status)
	}

	var result tokenResp
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", errors.New("خطا در بازخوانی اطلاعات آتی پی" + err.Error())
	}

	if result.Token != "" {
		return result.Token, nil
	}

	return "", errors.New(result.ErrorDesc + "---" + result.ErrorCode)
}

func (*Gate) SerializeToDB(resp interface{}) string {

	var result string
	switch v := resp.(type) {
	case map[string]string:
		result = strings.Join([]string{
			"State=" + v["state"],
			"Status=" + v["status"],
			"ReferenceNumber=" + v["referenceNumber"],
			"ReservationNumber=" + v["reservationNumber"],
			"TerminalID=" + v["terminalID"],
			"TraceNumber=" + v["traceNumber"],
		}, ";")
	case VerifyResp:
		result = strings.Join([]string{
			"Amount=" + strconv.Itoa(int(v.Amount)),
		}, ";")
	}

	return result
}

func (*Gate) CheckConfirmResp(resp map[string]string) (string, bool) {

	if resp["state"] != "OK" {
		return "", false
	}

	return resp["referenceNumber"], true
}

func (*Gate) CheckVerifyResp(resp interface{}) bool {

	return resp.(VerifyResp).Amount > 0
}

func (g *Gate) Verify(refNum string, resp interface{}) error {

	request := VerifyRequest{
		ReferenceNumber: refNum,
		TerminalID:      g.TerminalID,
		ApiKey:          g.ApiKey,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	var response *http.Response
	try := 0
	for ; try < 3; try++ {
		response, err = http.Post(consts.AtipayVerifyURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			sentrylog.Warning(errors.New("خطا در ارتباط با آتی پی " + err.Error()))
			continue
		}
		if response.StatusCode != http.StatusOK {
			sentrylog.Warning(errors.New("خطا در ارتباط با آتی پی " + response.Status))
			continue
		}
		break
	}
	defer io.SafeClose(response.Body)

	if try == 2 && (response == nil || response.StatusCode != http.StatusOK) {
		return errors.New("خطا در ارتباط با آتی پی " + response.Status)
	}

	var svr VerifyResp
	err = json.NewDecoder(response.Body).Decode(&svr)
	if err != nil {
		return errors.New("خطا در بازخوانی اطلاعات آتی پی" + err.Error())
	}

	reflect.ValueOf(resp).Elem().Set(reflect.ValueOf(svr))
	return nil
}
