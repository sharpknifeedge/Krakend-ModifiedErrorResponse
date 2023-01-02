package http

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/pg/internal/db"
	"gitlab.boomerangapp.ir/back/pg/internal/gateway/atipay"
	"gitlab.boomerangapp.ir/back/pg/internal/gateway/pec"
	"gitlab.boomerangapp.ir/back/pg/schema/entities"
	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/conv"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

func handleConfirm(w http.ResponseWriter, r *http.Request) {

	//get paymentID from url
	id, _ := strconv.ParseUint(mux.Vars(r)["paymentID"], 10, 32)

	//init db tx
	tx, err := db.DB.BeginTx(r.Context(), nil)
	if err != nil {
		sentrylog.Warning(err)
		tx.Rollback()
		return
	}

	//load payment
	payment, err := db.LoadPayment(r.Context(), tx, uint(id))
	if err != nil {
		sentrylog.Warning(err)
		tx.Rollback()
		return
	}

	//check payment status
	if payment.Status != consts.StatusIngress {
		sentrylog.Warning(errors.New("payment is not in progress"))
		tx.Rollback()
		return
	}

	//load and init psp gateway
	psp, _ := configs.GetGateConfig(payment.PSP)
	terminal := psp.Terminals[payment.GatewayID]

	//init gate interface via related psp
	var (
		gi         types.GateInterface
		verifyResp interface{}
	)
	switch payment.PSP {
	// case consts.PSPSadad:
	// 	gi = &sadad.Gate{ApiKey: terminal.ApiKey}
	// 	verifyResp = sadad.VerifyResp{}
	// case consts.PSPSep:
	// 	gi = &sep.Gate{MerchantID: psp.MerchantID}
	// 	verifyResp = sep.VerifyResp{}
	case consts.PSPAtipay:
		gi = &atipay.Gate{
			TerminalID: terminal.TerminalID,
			ApiKey:     terminal.ApiKey,
		}
		verifyResp = atipay.VerifyResp{}
	case consts.PSPPec:
		gi = &pec.Gate{ApiKey: terminal.ApiKey}
		verifyResp = pec.VerifyResp{}
	}

	//parsing confirm response form
	r.ParseForm()
	confirmResp := make(map[string]string)
	for k, v := range r.Form {
		confirmResp[k] = v[0]
	}

	//update payment record
	err = db.UpdateConfirm(r.Context(), tx, uint(id), gi.SerializeToDB(confirmResp))
	if err != nil {
		sentrylog.Warning(err)
		tx.Rollback()
		return
	}

	//check confirm response
	verifyCode, ok := gi.CheckConfirmResp(confirmResp)
	if !ok {
		handleConfirmErr(payment, tx, errors.New("error in confirmation"), "", w, r)
		return
	}

	//check verification with psp api
	err = gi.Verify(verifyCode, &verifyResp)
	if err != nil {
		handleConfirmErr(payment, tx, errors.New("error in calling verify"), "", w, r)
		return
	}

	//get verify values
	var (
		amount          uint
		refNo, cardMask string
	)
	switch v := verifyResp.(type) {
	case atipay.VerifyResp:
		amount = uint(v.Amount)
		refNo = confirmResp["referenceNumber"]
	case pec.VerifyResp:
		amt, _ := strconv.Atoi(strings.ReplaceAll(confirmResp["Amount"], ",", ""))
		amount = uint(amt)
		refNo = strconv.Itoa(int(v.Body.Resp.Result.RRN))
		cardMask = v.Body.Resp.Result.CardNumberMasked
		// case sadad.VerifyResp:
		// 	amount = v.Amount
		// 	refNo = v.RetrivalRefNo
		// case sep.VerifyResp:
		// 	amount = uint(v.Body.VerifyTransactionResponse.Result)
		// scr := confirmResp.(sep.ConfirmResp)
		// refNo = scr.RRN
	}

	//update payment record
	verifyrespLog := gi.SerializeToDB(verifyResp)
	err = db.UpdateVerify(r.Context(), tx, uint(id), amount, cardMask, verifyrespLog)
	if err != nil {
		sentrylog.Warning(errors.New(err.Error() + verifyrespLog))
		tx.Rollback()
		return
	}

	//check verify response
	if !gi.CheckVerifyResp(verifyResp) {
		handleConfirmErr(payment, tx, errors.New("error in verification"), refNo, w, r)
		return
	}

	//update status to verified
	err = db.UpdateStatus(r.Context(), tx, uint(id), consts.StatusVerified)
	if err != nil {
		sentrylog.Warning(errors.New(err.Error() + "status: verified"))
		tx.Rollback()
		return
	}

	//increase wallet balance
	err = db.IncreaseWalletBalance(r.Context(), tx, payment.WalletID, uint64(amount))
	if err != nil {
		sentrylog.Warning(err)
		tx.Rollback()
		return
	}

	//commit db tx
	err = tx.Commit()
	if err != nil {
		sentrylog.Warning(err)
		return
	}

	//write success page to user
	if payment.AutoRedirect {
		http.Redirect(w, r, payment.Callback, http.StatusTemporaryRedirect)
		return
	}
	file, err := ioutil.ReadFile(configs.Get().AppConfigs.HtmlSuccess)
	if err != nil {
		sentrylog.Warning(err)
		return
	}
	_, err = w.Write(conv.S2B(fmt.Sprintf(conv.B2S(file), amount, id, payment.Callback)))
	if err != nil {
		sentrylog.Warning(err)
	}
}

func handleConfirmErr(payment *entities.Payment, tx *sql.Tx, msg error, rrn string, w http.ResponseWriter, r *http.Request) {

	sentrylog.Warning(msg)

	err := db.UpdateStatus(r.Context(), tx, payment.ID, consts.StatusFailed)
	if err != nil {
		sentrylog.Warning(err)
	}
	tx.Commit()

	if payment.AutoRedirect {
		http.Redirect(w, r, payment.Callback, http.StatusTemporaryRedirect)
		return
	}

	file, err := ioutil.ReadFile(configs.Get().AppConfigs.Htmlfailed)
	if err != nil {
		sentrylog.Warning(err)
		return
	}
	_, err = w.Write(conv.S2B(fmt.Sprintf(conv.B2S(file), rrn, payment.ID, payment.Callback)))
	if err != nil {
		sentrylog.Warning(err)
	}
}
