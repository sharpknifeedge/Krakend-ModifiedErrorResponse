package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/pg/internal/db"
	"gitlab.boomerangapp.ir/back/pg/internal/gateway/atipay"
	"gitlab.boomerangapp.ir/back/pg/internal/gateway/pec"
	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	utilsConsts "gitlab.boomerangapp.ir/back/utils/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/http_helper"
	"gitlab.boomerangapp.ir/back/utils/pkg/response"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
	"gitlab.boomerangapp.ir/back/utils/pkg/user_grpc_client"
	"gitlab.boomerangapp.ir/back/utils/pkg/validation"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

	//init response handler
	mr := http_helper.MuxResp{
		Writer: w,
		Hub:    sentrylog.GetHubFromCtx(r.Context()),
	}

	//get user|producer id from header
	role := r.Header.Get(utilsConsts.HTTPROleHeader)
	if role != utilsConsts.RoleProducer && role != utilsConsts.RoleUser {
		mr.HandleErr(nil, response.AccessDenied)
		return
	}
	userID, _ := strconv.ParseUint(r.Header.Get(utilsConsts.HTTPUserID), 10, 32)
	ip := r.Header.Get(utilsConsts.HTTPXForwardedFor)

	//get payment request from body
	var pr types.PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		mr.HandleErr(err, response.UnprocessableEntity)
		return
	}

	//check payment rules
	mD, err := db.LoadMinDeposit(r.Context())
	if err != nil {
		mr.HandleErr(err, response.InternalError)
		return
	}
	if pr.Amount < mD {
		mr.HandleServiceErr(consts.ErrCodeLessThanMinAmount, consts.Messages)
		return
	}
	maxBalance, maxDeposit, err := db.LoadMaxRules(r.Context(), uint(userID))
	if err != nil {
		mr.HandleErr(err, response.InternalError)
		return
	}
	if maxBalance != 0 {
		//check max wallet balance
		walletBalance, err := db.LoadWalletBalance(r.Context(), uint(userID))
		if err != nil {
			mr.HandleErr(err, response.InternalError)
			return
		}
		if walletBalance+int64(pr.Amount) > int64(maxBalance) {
			mr.HandleServiceErr(consts.ErrCodeMoreThanMaxBalance, consts.Messages)
			return
		}
	}
	if maxDeposit != 0 {
		//check max daily deposit
		depositSum, err := db.LoadDailyDepositSum(r.Context(), uint(userID))
		if err != nil {
			mr.HandleErr(err, response.InternalError)
			return
		}
		if depositSum+uint64(pr.Amount) > uint64(maxDeposit) {
			mr.HandleServiceErr(consts.ErrCodeMoreThanDailyDeposit, consts.Messages)
			return
		}
	}

	//call user RPC
	user, err := user_grpc_client.GetUserById(configs.Get().AppConfigs.GrpcUserAddr, uint(userID))
	if err != nil {
		mr.HandleServiceErr(consts.ErrCodeInvalidUserID, consts.Messages)
		return
	}

	//load psp
	psp, ok := configs.GetGateConfig(pr.PSP)
	if !ok {
		mr.HandleServiceErr(consts.ErrCodeInvalidPSP, consts.Messages)
		return
	}

	//check and init terminal
	if int(pr.GatewayID) > len(psp.Terminals) {
		mr.HandleServiceErr(consts.ErrCodeInvalidTerminal, consts.Messages)
		return
	}
	terminal := psp.Terminals[pr.GatewayID]

	//init gateway interface
	var gi types.GateInterface
	switch pr.PSP {
	// case consts.PSPSadad:
	// 	gi = &sadad.Gate{
	// 		MerchantID: psp.MerchantID,
	// 		TerminalID: terminal.TerminalID,
	// 		ApiKey:     terminal.ApiKey,
	// 	}
	// case consts.PSPSep:
	// 	gi = &sep.Gate{TerminalID: terminal.TerminalID}
	case consts.PSPAtipay:
		gi = &atipay.Gate{ApiKey: terminal.ApiKey}
	case consts.PSPPec:
		gi = &pec.Gate{ApiKey: terminal.ApiKey}
	}

	//init db record
	pr.ID, err = db.CreatePayment(r.Context(), pr, role, ip, uint(userID))
	if err != nil {
		vErr := validation.ValidateByErr(err)
		if vErr != nil {
			mr.HandleValidationErr(vErr)
		} else {
			mr.HandleErr(err, response.InternalError)
		}
		return
	}

	//get token from psp api
	token, err := gi.GetToken(pr, user.Mobile)
	if err != nil {
		mr.HandleServiceErr(consts.ErrCodeUnableToGetToken, consts.Messages)
		return
	}

	//update db with token
	err = db.UpdateToken(r.Context(), pr.ID, token)
	if err != nil {
		vErr := validation.ValidateByErr(err)
		if vErr != nil {
			mr.HandleValidationErr(vErr)
		} else {
			mr.HandleErr(err, response.InternalError)
		}
		return
	}

	//send response
	mr.SendResp(struct {
		ID  uint   `json:"id"`
		URL string `json:"url"`
	}{
		pr.ID,
		gi.MakeRedirectURL(token),
	}, nil, true)
}
